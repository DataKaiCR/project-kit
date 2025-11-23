package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/datakaicr/pk/pkg/config"
	"github.com/datakaicr/pk/pkg/paths"
)

// AccessRecord tracks when a project was last accessed
type AccessRecord struct {
	ProjectID    string    `json:"project_id"`
	ProjectPath  string    `json:"project_path"`
	LastAccessed time.Time `json:"last_accessed"`
}

// GetAccessFile returns the path to the access tracking file
func GetAccessFile() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	cacheDir := filepath.Join(homeDir, ".cache", "pk")
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return "", err
	}

	return filepath.Join(cacheDir, "access.json"), nil
}

// LoadAccessRecords reads the access tracking file and validates paths
// Automatically heals stale paths by searching for projects
func LoadAccessRecords() (map[string]AccessRecord, error) {
	accessFile, err := GetAccessFile()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(accessFile)
	if err != nil {
		if os.IsNotExist(err) {
			return make(map[string]AccessRecord), nil
		}
		return nil, err
	}

	var records map[string]AccessRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, err
	}

	// Validate and heal paths
	healed, err := validateAndHealAccessRecords(records)
	if err != nil {
		// If validation fails, return original records
		return records, nil
	}

	// If any paths were healed, save updated records
	if healed {
		SaveAccessRecords(records)
	}

	return records, nil
}

// validateAndHealAccessRecords checks if access record paths exist and updates them if stale
// Returns true if any paths were healed
func validateAndHealAccessRecords(records map[string]AccessRecord) (bool, error) {
	resolver, err := paths.NewResolver()
	if err != nil {
		return false, err
	}

	healed := false
	for projectID, record := range records {
		newPath, wasHealed, err := resolver.ValidatePath(record.ProjectID, record.ProjectPath)
		if err != nil {
			// Project not found, keep original path
			continue
		}

		if wasHealed {
			record.ProjectPath = newPath
			records[projectID] = record
			healed = true
		}
	}

	return healed, nil
}

// SaveAccessRecords writes the access tracking file
func SaveAccessRecords(records map[string]AccessRecord) error {
	accessFile, err := GetAccessFile()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(records, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(accessFile, data, 0644)
}

// RecordAccess marks a project as accessed now
func RecordAccess(projectID, projectPath string) error {
	records, err := LoadAccessRecords()
	if err != nil {
		return err
	}

	records[projectID] = AccessRecord{
		ProjectID:    projectID,
		ProjectPath:  projectPath,
		LastAccessed: time.Now(),
	}

	return SaveAccessRecords(records)
}

// GetRecentProjects returns projects sorted by access time (most recent first)
func GetRecentProjects(limit int) ([]*config.Project, error) {
	// Load access records
	records, err := LoadAccessRecords()
	if err != nil {
		return nil, err
	}

	// Load all projects
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	projects, err := FindProjectsCached(
		filepath.Join(homeDir, "projects"),
		filepath.Join(homeDir, "scratch"),
	)
	if err != nil {
		return nil, err
	}

	// Sort projects by access time
	sort.Slice(projects, func(i, j int) bool {
		accessI, okI := records[projects[i].ProjectInfo.ID]
		accessJ, okJ := records[projects[j].ProjectInfo.ID]

		// Projects never accessed go to the end
		if !okI && !okJ {
			return false
		}
		if !okI {
			return false
		}
		if !okJ {
			return true
		}

		return accessI.LastAccessed.After(accessJ.LastAccessed)
	})

	// Apply limit
	if limit > 0 && limit < len(projects) {
		projects = projects[:limit]
	}

	return projects, nil
}
