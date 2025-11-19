package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/datakaicr/pk/pkg/config"
	"github.com/spf13/cobra"
)

var archiveCmd = &cobra.Command{
	Use:   "archive <name>",
	Short: "Archive a project",
	Long: `Move a project to the archive directory and update its status.

This will:
  1. Move the project from ~/projects to ~/archive
  2. Update status to "archived" in .project.toml
  3. Set completion date to today
  4. Auto-sync shell aliases (if enabled)

Example:
  pk archive old-project
  pk archive keplr-data-model`,
	Args:              cobra.ExactArgs(1),
	Run:               runArchive,
	ValidArgsFunction: validProjectNames,
}

var archiveAutoSync bool

func init() {
	rootCmd.AddCommand(archiveCmd)
	archiveCmd.Flags().BoolVar(&archiveAutoSync, "sync", true, "Auto-sync aliases after archiving")
}

func runArchive(cmd *cobra.Command, args []string) {
	projectName := strings.ToLower(args[0])

	homeDir, _ := os.UserHomeDir()
	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")

	// Find project in projects directory
	projects, err := config.FindProjects(projectsDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error finding projects: %v\n", err)
		os.Exit(1)
	}

	var found *config.Project
	for _, p := range projects {
		if strings.ToLower(p.ProjectInfo.ID) == projectName ||
			strings.ToLower(p.ProjectInfo.Name) == projectName {
			found = p
			break
		}
	}

	if found == nil {
		fmt.Fprintf(os.Stderr, "Project '%s' not found in ~/projects\n", projectName)
		fmt.Fprintf(os.Stderr, "Hint: Use 'pk list active' to see available projects\n")
		os.Exit(1)
	}

	// Check if already exists in archive
	destPath := filepath.Join(archiveDir, filepath.Base(found.Path))
	if _, err := os.Stat(destPath); !os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Project already exists in archive: %s\n", destPath)
		os.Exit(1)
	}

	// Ensure archive directory exists
	if err := os.MkdirAll(archiveDir, 0755); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create archive directory: %v\n", err)
		os.Exit(1)
	}

	// Move project
	fmt.Printf("Moving project: %s\n", found.ProjectInfo.Name)
	fmt.Printf("  From: %s\n", found.Path)
	fmt.Printf("  To:   %s\n", destPath)

	if err := os.Rename(found.Path, destPath); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to move project: %v\n", err)
		os.Exit(1)
	}

	// Update .project.toml
	tomlPath := filepath.Join(destPath, ".project.toml")
	if err := updateProjectToml(tomlPath); err != nil {
		fmt.Fprintf(os.Stderr, "Warning: Failed to update .project.toml: %v\n", err)
	} else {
		fmt.Printf("\n\033[32m✓\033[0m Archived successfully\n")
		fmt.Printf("  Status: \033[33marchived\033[0m\n")
		fmt.Printf("  Location: %s\n", destPath)
	}

	// Auto-sync aliases
	if archiveAutoSync {
		fmt.Printf("\nSyncing aliases...\n")
		runSync(cmd, []string{})
	}
}

func updateProjectToml(path string) error {
	// Read current TOML
	var project config.Project
	if _, err := toml.DecodeFile(path, &project); err != nil {
		return err
	}

	// Update status and completion date
	project.ProjectInfo.Status = "archived"
	project.Dates.Completed = time.Now().Format("2006-01-02")

	// Write back to file
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	encoder := toml.NewEncoder(f)
	return encoder.Encode(&project)
}
