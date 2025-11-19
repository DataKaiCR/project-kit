package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/datakaicr/pk/pkg/config"
	"github.com/datakaicr/pk/pkg/session"
	"github.com/spf13/cobra"
)

var (
	deleteKeepGit bool
	deleteForce   bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <name>",
	Short: "Delete a project",
	Long: `Remove a project directory and all its contents.

This will:
  1. Validate project exists
  2. Check for active tmux session and optionally kill it
  3. Optionally archive git history (--keep-git)
  4. Remove entire project directory
  5. Auto-sync shell aliases

WARNING: This operation is permanent. Data will be deleted.

Example:
  pk delete old-project
  pk delete legacy-project --force         # Skip confirmation, auto-kill session
  pk delete archived-proj --keep-git       # Save git history first`,
	Args:              cobra.ExactArgs(1),
	Run:               runDelete,
	ValidArgsFunction: validProjectNames,
}

func init() {
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().BoolVar(&deleteKeepGit, "keep-git", false,
		"Archive git history before deletion")
	deleteCmd.Flags().BoolVar(&deleteForce, "force", false,
		"Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) {
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

	// Check for active tmux session
	sessionName := session.SanitizeSessionName(found.ProjectInfo.ID)
	hasSession := session.SessionExists(sessionName)

	// Show confirmation prompt
	if !deleteForce {
		fmt.Printf("\033[33mWARNING: This will permanently delete the project.\033[0m\n\n")
		fmt.Printf("Project:  %s\n", found.ProjectInfo.Name)
		fmt.Printf("Location: %s\n", found.Path)
		fmt.Printf("Status:   %s\n", found.ProjectInfo.Status)
		if hasSession {
			fmt.Printf("Tmux:     \033[33m● Active session found\033[0m\n")
		}
		fmt.Println()

		fmt.Print("Continue? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" {
			fmt.Println("Cancelled")
			return
		}
	}

	// Kill tmux session if it exists
	if hasSession {
		if !deleteForce {
			fmt.Print("\nKill active tmux session? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if strings.ToLower(response) == "y" {
				if err := session.KillSession(sessionName); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: Failed to kill tmux session: %v\n", err)
				} else {
					fmt.Printf("\033[32m✓\033[0m Tmux session killed\n")
				}
			} else {
				fmt.Println("Tmux session will remain active")
			}
		} else {
			// Force flag: auto-kill session
			if err := session.KillSession(sessionName); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to kill tmux session: %v\n", err)
			} else {
				fmt.Printf("\033[32m✓\033[0m Tmux session killed\n")
			}
		}
	}

	// Archive git history if requested
	if deleteKeepGit {
		gitDir := filepath.Join(found.Path, ".git")
		if _, err := os.Stat(gitDir); err == nil {
			archiveName := filepath.Base(found.Path) + ".git-archive.tar.gz"
			archivePath := filepath.Join(filepath.Dir(found.Path), archiveName)

			fmt.Printf("Archiving git history to: %s\n", archivePath)

			tarCmd := exec.Command("tar", "czf", archivePath, "-C", found.Path, ".git")
			if err := tarCmd.Run(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: Failed to archive git history: %v\n", err)
				fmt.Print("Continue with deletion? (y/N): ")

				var response string
				fmt.Scanln(&response)

				if strings.ToLower(response) != "y" {
					fmt.Println("Cancelled")
					return
				}
			} else {
				fmt.Printf("\033[32m✓\033[0m Git history archived\n")
			}
		} else {
			fmt.Println("No git repository found, skipping archive")
		}
	}

	// Delete project directory
	if err := os.RemoveAll(found.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to delete project: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\033[32m✓\033[0m Deleted: %s\n", found.Path)

	// Sync aliases
	fmt.Println("Syncing aliases...")
	runSync(cmd, []string{})

	fmt.Printf("\n\033[32m✓\033[0m Project '%s' deleted successfully\n", found.ProjectInfo.Name)
}
