package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/datakaicr/pk/pkg/config"
	"github.com/spf13/cobra"
)

var editCmd = &cobra.Command{
	Use:   "edit <name>",
	Short: "Edit project metadata",
	Long: `Open the project's .project.toml file in your default editor.

The editor is determined by (in order):
  1. $EDITOR environment variable
  2. vim
  3. nano

After editing, the TOML is validated. If the project ID changed,
aliases will be regenerated automatically.

Example:
  pk edit dojo
  pk edit my-project`,
	Args:              cobra.ExactArgs(1),
	Run:               runEdit,
	ValidArgsFunction: validProjectNames,
}

func init() {
	rootCmd.AddCommand(editCmd)
}

func runEdit(cmd *cobra.Command, args []string) {
	projectName := strings.ToLower(args[0])

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
		if strings.ToLower(p.ProjectInfo.ID) == projectName ||
			strings.ToLower(p.ProjectInfo.Name) == projectName {
			found = p
			break
		}
	}

	if found == nil {
		fmt.Fprintf(os.Stderr, "Error: Project '%s' not found\n", args[0])
		fmt.Fprintf(os.Stderr, "\nUse 'pk list' to see all projects.\n")
		os.Exit(1)
	}

	tomlPath := filepath.Join(found.Path, ".project.toml")

	// Store original ID to detect changes
	originalID := found.ProjectInfo.ID

	// Determine editor
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vim"
		// Check if vim exists, fallback to nano
		if _, err := exec.LookPath("vim"); err != nil {
			editor = "nano"
		}
	}

	fmt.Printf("Opening %s in %s...\n", tomlPath, editor)

	// Open editor
	editorCmd := exec.Command(editor, tomlPath)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Editor failed: %v\n", err)
		os.Exit(1)
	}

	// Validate TOML
	var project config.Project
	if _, err := toml.DecodeFile(tomlPath, &project); err != nil {
		fmt.Fprintf(os.Stderr, "\n\033[33mWarning: Invalid TOML syntax:\033[0m %v\n", err)
		fmt.Fprintf(os.Stderr, "Please fix the file and run 'pk sync' when ready.\n")
		os.Exit(1)
	}

	fmt.Printf("\n\033[32m✓\033[0m Metadata updated successfully\n")

	// Check if ID changed
	if project.ProjectInfo.ID != originalID {
		fmt.Printf("\nProject ID changed: %s → %s\n", originalID, project.ProjectInfo.ID)
		fmt.Println("Syncing aliases...")
		runSync(cmd, []string{})
	}

	fmt.Printf("\nProject: %s\n", project.ProjectInfo.Name)
	fmt.Printf("Status:  %s\n", project.ProjectInfo.Status)
	fmt.Printf("Type:    %s\n", project.ProjectInfo.Type)
}
