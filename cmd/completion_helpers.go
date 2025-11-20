package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/datakaicr/pk/pkg/cache"
	"github.com/spf13/cobra"
)

// validProjectNames returns list of project names/IDs for completion
func validProjectNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")
	scriptoriumDir := filepath.Join(homeDir, "scriptorium")

	// Use cached projects if available
	projects, err := cache.FindProjectsCached(projectsDir, archiveDir, scriptoriumDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, p := range projects {
		// Add project ID
		if strings.HasPrefix(p.ProjectInfo.ID, toComplete) {
			names = append(names, p.ProjectInfo.ID)
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

// validScratchNames returns list of scratch project names for completion
func validScratchNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	scratchDir := filepath.Join(homeDir, "scratch")

	// Check if scratch directory exists
	if _, err := os.Stat(scratchDir); os.IsNotExist(err) {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Read directories
	entries, err := os.ReadDir(scratchDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, entry := range entries {
		if entry.IsDir() && strings.HasPrefix(entry.Name(), toComplete) {
			names = append(names, entry.Name())
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

// validAllProjectNames returns both regular projects and scratch projects
func validAllProjectNames(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")
	scriptoriumDir := filepath.Join(homeDir, "scriptorium")
	scratchDir := filepath.Join(homeDir, "scratch")

	// Get regular projects
	projects, err := cache.FindProjectsCached(projectsDir, archiveDir, scriptoriumDir)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	var names []string
	for _, p := range projects {
		if strings.HasPrefix(p.ProjectInfo.ID, toComplete) {
			names = append(names, p.ProjectInfo.ID)
		}
	}

	// Get scratch projects
	if _, err := os.Stat(scratchDir); err == nil {
		entries, err := os.ReadDir(scratchDir)
		if err == nil {
			for _, entry := range entries {
				if entry.IsDir() && strings.HasPrefix(entry.Name(), toComplete) {
					names = append(names, entry.Name())
				}
			}
		}
	}

	return names, cobra.ShellCompDirectiveNoFileComp
}

// validListFilters returns valid filter options for pk list
func validListFilters(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	filters := []string{"active", "archived", "datakai", "westmonroe", "product", "client"}
	var matches []string
	for _, f := range filters {
		if strings.HasPrefix(f, toComplete) {
			matches = append(matches, f)
		}
	}
	return matches, cobra.ShellCompDirectiveNoFileComp
}
