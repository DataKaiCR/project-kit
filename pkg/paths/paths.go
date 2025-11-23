package paths

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Config holds user-configurable paths
type Config struct {
	Paths struct {
		Projects    string `toml:"projects"`
		Archive     string `toml:"archive"`
		Scratch     string `toml:"scratch"`
		Scriptorium string `toml:"scriptorium"`
	} `toml:"paths"`
}

// Resolver handles path resolution with config and defaults
type Resolver struct {
	homeDir     string
	config      *Config
	projects    string
	archive     string
	scratch     string
	scriptorium string
}

// NewResolver creates a new path resolver
// Loads config from ~/.config/pk/config.toml if it exists, otherwise uses defaults
func NewResolver() (*Resolver, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	r := &Resolver{
		homeDir: homeDir,
	}

	// Try to load config
	configPath := filepath.Join(homeDir, ".config", "pk", "config.toml")
	if _, err := os.Stat(configPath); err == nil {
		// Config exists, load it
		var cfg Config
		if _, err := toml.DecodeFile(configPath, &cfg); err != nil {
			// Config malformed, warn but continue with defaults
			fmt.Fprintf(os.Stderr, "Warning: Failed to parse config %s: %v\n", configPath, err)
		} else {
			r.config = &cfg
		}
	}

	// Set paths (use config if available, otherwise defaults)
	r.projects = r.resolvePath("projects", filepath.Join(homeDir, "projects"))
	r.archive = r.resolvePath("archive", filepath.Join(homeDir, "archive"))
	r.scratch = r.resolvePath("scratch", filepath.Join(homeDir, "scratch"))
	r.scriptorium = r.resolvePath("scriptorium", filepath.Join(homeDir, "scriptorium"))

	return r, nil
}

// resolvePath returns config path if set, otherwise returns defaultPath
// Expands ~ to home directory
func (r *Resolver) resolvePath(name, defaultPath string) string {
	if r.config == nil {
		return defaultPath
	}

	var configured string
	switch name {
	case "projects":
		configured = r.config.Paths.Projects
	case "archive":
		configured = r.config.Paths.Archive
	case "scratch":
		configured = r.config.Paths.Scratch
	case "scriptorium":
		configured = r.config.Paths.Scriptorium
	}

	if configured == "" {
		return defaultPath
	}

	// Expand ~ to home directory
	if configured[0] == '~' {
		configured = filepath.Join(r.homeDir, configured[1:])
	}

	return configured
}

// Projects returns the projects directory path
func (r *Resolver) Projects() string {
	return r.projects
}

// Archive returns the archive directory path
func (r *Resolver) Archive() string {
	return r.archive
}

// Scratch returns the scratch directory path
func (r *Resolver) Scratch() string {
	return r.scratch
}

// Scriptorium returns the scriptorium directory path
func (r *Resolver) Scriptorium() string {
	return r.scriptorium
}

// AllRoots returns all root directories
func (r *Resolver) AllRoots() []string {
	return []string{
		r.projects,
		r.archive,
		r.scriptorium,
	}
}

// FindProject searches for a project by ID across all root directories
// Returns the full path if found, empty string if not found
func (r *Resolver) FindProject(projectID string) (string, error) {
	// Check each root directory
	for _, root := range r.AllRoots() {
		// Check if root exists
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}

		// Read directory entries
		entries, err := os.ReadDir(root)
		if err != nil {
			continue
		}

		// Look for matching project
		for _, entry := range entries {
			if !entry.IsDir() {
				continue
			}

			projectPath := filepath.Join(root, entry.Name())

			// Check if directory matches project ID
			if entry.Name() == projectID {
				return projectPath, nil
			}

			// Also check .project.toml for ID match
			tomlPath := filepath.Join(projectPath, ".project.toml")
			if _, err := os.Stat(tomlPath); err == nil {
				// TODO: Parse .project.toml and check ID
				// For now, rely on directory name match
			}
		}
	}

	// Also check scratch directory (different structure)
	scratchPath := filepath.Join(r.scratch, projectID)
	if _, err := os.Stat(scratchPath); err == nil {
		return scratchPath, nil
	}

	return "", fmt.Errorf("project %s not found", projectID)
}

// ValidatePath checks if a path exists, and if not, attempts to find the project by ID
// Returns the validated path (original if valid, or new path if found)
// Returns error if project cannot be found
func (r *Resolver) ValidatePath(projectID, cachedPath string) (string, bool, error) {
	// Check if cached path still exists
	if _, err := os.Stat(cachedPath); err == nil {
		return cachedPath, false, nil // Path is valid, no healing needed
	}

	// Path is stale, try to find project
	newPath, err := r.FindProject(projectID)
	if err != nil {
		return "", false, err
	}

	return newPath, true, nil // Path was healed
}

// Default returns a resolver with default settings (no config)
// This is useful for testing or when config loading fails
func Default() (*Resolver, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}

	return &Resolver{
		homeDir:     homeDir,
		projects:    filepath.Join(homeDir, "projects"),
		archive:     filepath.Join(homeDir, "archive"),
		scratch:     filepath.Join(homeDir, "scratch"),
		scriptorium: filepath.Join(homeDir, "scriptorium"),
	}, nil
}
