package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pk",
	Short: "Project Kit - Manage projects with .project.toml metadata",
	Long: `PK (Project Kit) is a cross-platform CLI tool for managing projects.

It discovers projects by .project.toml files, helps with archiving,
generates shell aliases, and provides project statistics.

Commands:
  pk scratch new <name>    # Create scratch project for experimentation
  pk new <name>            # Create a new project
  pk list [filter]     # List all projects (active, archived, datakai, etc.)
  pk show <name>       # Show detailed project information
  pk edit <name>       # Edit project metadata
  pk rename <old> <new>  # Rename a project
  pk promote <path>    # Convert directory into a project
  pk archive <name>    # Archive a project (move to ~/archive)
  pk delete <name>     # Delete a project permanently
  pk sync              # Generate shell aliases for all projects

Workflow:
  pk scratch new prototype      # Quick experimentation
  pk promote prototype          # Ready for real work
  pk edit prototype             # Update metadata
  pk rename prototype app       # Better name

Examples:
  pk scratch new api-test       # Create in ~/scratch
  pk scratch list               # List all scratch projects
  pk scratch delete old-test    # Remove scratch project
  pk new my-project         # Create in ~/projects
  pk list active            # List active projects
  pk show dojo              # Show project details
  pk edit dojo              # Edit metadata in $EDITOR
  pk rename old-name new    # Rename project
  pk promote api-test       # Promote scratch to project
  pk archive old-proj       # Archive a project
  pk delete test --force    # Delete without confirmation`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Global flags (available to all commands)
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.pk.yaml)")

	// Local flags (only for this command)
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
