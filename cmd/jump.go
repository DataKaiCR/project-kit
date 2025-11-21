package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/datakaicr/pk/pkg/cache"
	"github.com/datakaicr/pk/pkg/config"
	"github.com/datakaicr/pk/pkg/context"
	"github.com/datakaicr/pk/pkg/session"
	"github.com/spf13/cobra"
)

var jumpCmd = &cobra.Command{
	Use:   "jump <slot>",
	Short: "Jump to a pinned project by slot number",
	Long: `Jump to a pinned project by its slot number (1-5).

This command opens the pinned project in a tmux session, creating one if needed.
Projects must first be pinned with 'pk pin add <project> <slot>'.

Designed for keyboard shortcuts in tmux:
  bind-key g switch-client -T jump
  bind-key -T jump 1 run-shell "pk jump 1"
  bind-key -T jump 2 run-shell "pk jump 2"
  ...

Then use: Ctrl+b g 1, Ctrl+b g 2, etc.

Examples:
  pk jump 1     # Jump to project in slot 1
  pk jump 2     # Jump to project in slot 2`,
	Args:              cobra.ExactArgs(1),
	Run:               runJump,
	ValidArgsFunction: validJumpArgs,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return session.CheckTmux()
	},
}

func init() {
	rootCmd.AddCommand(jumpCmd)
}

func runJump(cmd *cobra.Command, args []string) {
	slotStr := args[0]

	// Parse slot number
	slot, err := strconv.Atoi(slotStr)
	if err != nil || slot < 1 || slot > 5 {
		fmt.Fprintf(os.Stderr, "Error: Slot must be a number between 1 and 5\n")
		os.Exit(1)
	}

	// Get the pin
	pin, err := cache.GetPin(slot)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		fmt.Fprintf(os.Stderr, "\nPin a project first with:\n")
		fmt.Fprintf(os.Stderr, "  pk pin add <project> %d\n", slot)
		os.Exit(1)
	}

	// Load project configuration
	projectToml := pin.ProjectPath + "/.project.toml"
	project, err := config.LoadProject(projectToml)
	if err != nil {
		// Create a basic project structure if .project.toml doesn't exist
		project = &config.Project{
			Path: pin.ProjectPath,
		}
		project.ProjectInfo.ID = pin.ProjectID
		project.ProjectInfo.Name = pin.ProjectID
		project.ProjectInfo.Status = "active"
	}

	// Record access
	cache.RecordAccess(pin.ProjectID, pin.ProjectPath)

	// Switch context if configured
	context.Switch(project)

	// Create or switch to session
	if err := session.CreateSession(project); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to create/switch session: %v\n", err)
		os.Exit(1)
	}
}

// validJumpArgs provides shell completion for jump command
func validJumpArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// Show available slots with their pinned projects
		pins, err := cache.ListPins()
		if err != nil {
			return []string{"1", "2", "3", "4", "5"}, cobra.ShellCompDirectiveNoFileComp
		}

		var completions []string
		pinMap := make(map[int]string)
		for _, pin := range pins {
			pinMap[pin.Slot] = pin.ProjectID
		}

		for i := 1; i <= 5; i++ {
			if projectID, exists := pinMap[i]; exists {
				completions = append(completions, fmt.Sprintf("%d\t%s", i, projectID))
			} else {
				completions = append(completions, fmt.Sprintf("%d\t(empty)", i))
			}
		}

		return completions, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}
