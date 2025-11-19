package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/datakaicr/pk/pkg/config"
	"github.com/spf13/cobra"
)

var renameCmd = &cobra.Command{
	Use:   "rename <old-name> <new-name>",
	Short: "Rename a project",
	Long: `Rename a project directory and update all metadata.

This will:
  1. Validate both old and new names
  2. Rename the project directory
  3. Update .project.toml (name and ID)
  4. Auto-sync shell aliases

Example:
  pk rename old-name new-name
  pk rename prototype awesome-product`,
	Args: cobra.ExactArgs(2),
	Run:  runRename,
}

func init() {
	rootCmd.AddCommand(renameCmd)
}

func runRename(cmd *cobra.Command, args []string) {
	oldName := strings.ToLower(args[0])
	newName := args[1]

	// Validate new name
	if strings.ContainsAny(newName, "/\\:*?\"<>|") {
		fmt.Fprintf(os.Stderr, "Error: Invalid project name. Avoid special characters.\n")
		os.Exit(1)
	}

	// Find project
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Could not determine home directory: %v\n", err)
		os.Exit(1)
	}

	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")

	projects, err := config.FindProjects(projectsDir, archiveDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to find projects: %v\n", err)
		os.Exit(1)
	}

	var found *config.Project
	for _, p := range projects {
		if strings.ToLower(p.ProjectInfo.ID) == oldName ||
			strings.ToLower(p.ProjectInfo.Name) == oldName {
			found = p
			break
		}
	}

	if found == nil {
		fmt.Fprintf(os.Stderr, "Error: Project '%s' not found\n", args[0])
		fmt.Fprintf(os.Stderr, "\nUse 'pk list' to see all projects.\n")
		os.Exit(1)
	}

	// Determine new path
	parentDir := filepath.Dir(found.Path)
	newPath := filepath.Join(parentDir, newName)

	// Check if new name already exists
	if _, err := os.Stat(newPath); err == nil {
		fmt.Fprintf(os.Stderr, "Error: A project with name '%s' already exists at %s\n", newName, newPath)
		os.Exit(1)
	}

	fmt.Printf("Renaming project: %s → %s\n", found.ProjectInfo.Name, newName)
	fmt.Printf("Location: %s → %s\n", found.Path, newPath)

	// Rename directory
	if err := os.Rename(found.Path, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to rename directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m✓\033[0m Directory renamed\n")

	// Update .project.toml
	tomlPath := filepath.Join(newPath, ".project.toml")
	if err := updateProjectTomlRename(tomlPath, newName, newPath); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to update .project.toml: %v\n", err)
		fmt.Fprintf(os.Stderr, "Directory was renamed but metadata update failed.\n")
		os.Exit(1)
	}

	fmt.Printf("\033[32m✓\033[0m Metadata updated\n")

	// Sync aliases
	fmt.Println("Syncing aliases...")
	runSync(cmd, []string{})

	fmt.Printf("\n\033[32m✓\033[0m Project renamed successfully!\n")
	fmt.Printf("\nNew alias:\n")
	fmt.Printf("  %s    # Jump to project (after reloading shell)\n", newName)
}

func updateProjectTomlRename(path, newName, newPath string) error {
	// Read current TOML
	var project config.Project
	if _, err := toml.DecodeFile(path, &project); err != nil {
		return err
	}

	// Update fields
	project.Path = newPath
	project.ProjectInfo.Name = newName
	project.ProjectInfo.ID = newName

	// Write back
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	// Write header comment
	fmt.Fprintln(f, "# Project Metadata")
	fmt.Fprintln(f, "")

	encoder := toml.NewEncoder(f)
	return encoder.Encode(&project)
}
