package hooks

import (
	"os"
	"path/filepath"

	"github.com/datakaicr/pk/pkg/cache"
)

// InvalidateCache triggers a cache rebuild after project modifications
func InvalidateCache() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return
	}

	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")
	scriptoriumDir := filepath.Join(homeDir, "scriptorium")

	// Rebuild cache in background
	cache.RebuildCacheAsync(projectsDir, archiveDir, scriptoriumDir)
}
