package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/datakaicr/pk/pkg/cache"
	"github.com/datakaicr/pk/pkg/config"
	"github.com/spf13/cobra"
)

var pinCmd = &cobra.Command{
	Use:   "pin",
	Short: "Manage pinned projects for quick jumping",
	Long: `Pin projects to slots (1-5) for instant access with 'pk jump'.

Pinned projects can be quickly accessed with keyboard shortcuts:
  Ctrl+b g 1  # Jump to pin slot 1
  Ctrl+b g 2  # Jump to pin slot 2
  etc.

Subcommands:
  pk pin add <project> <slot>   # Pin a project to a slot
  pk pin remove <slot|project>  # Remove a pin
  pk pin list                   # Show all pins
  pk pin clear                  # Remove all pins`,
}

var pinAddCmd = &cobra.Command{
	Use:   "add <project> <slot>",
	Short: "Pin a project to a slot (1-5)",
	Long: `Pin a project to a numbered slot for quick access.

Slots are numbered 1-5 and can be accessed via:
  pk jump 1
  pk jump 2
  etc.

Or with tmux keybindings:
  Ctrl+b g 1
  Ctrl+b g 2

Examples:
  pk pin add pk 1          # Pin 'pk' to slot 1
  pk pin add dkos 2        # Pin 'dkos' to slot 2
  pk pin add conduit 3     # Pin 'conduit' to slot 3`,
	Args:              cobra.ExactArgs(2),
	Run:               runPinAdd,
	ValidArgsFunction: validPinAddArgs,
}

var pinRemoveCmd = &cobra.Command{
	Use:   "remove <slot|project>",
	Short: "Remove a pin by slot number or project name",
	Long: `Remove a pinned project by slot number or project name.

Examples:
  pk pin remove 1       # Remove pin in slot 1
  pk pin remove pk      # Remove pin for project 'pk'`,
	Args: cobra.ExactArgs(1),
	Run:  runPinRemove,
}

var pinListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all pinned projects",
	Long: `Display all currently pinned projects with their slot numbers.

Shows which projects are pinned to which slots for quick reference.`,
	Run: runPinList,
}

var pinClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Remove all pins",
	Long: `Clear all pinned projects at once.

This will remove all pins from all slots. Use with caution.`,
	Run: runPinClear,
}

func init() {
	rootCmd.AddCommand(pinCmd)
	pinCmd.AddCommand(pinAddCmd)
	pinCmd.AddCommand(pinRemoveCmd)
	pinCmd.AddCommand(pinListCmd)
	pinCmd.AddCommand(pinClearCmd)
}

func runPinAdd(cmd *cobra.Command, args []string) {
	projectName := strings.ToLower(args[0])
	slotStr := args[1]

	// Parse slot number
	slot, err := strconv.Atoi(slotStr)
	if err != nil || slot < 1 || slot > 5 {
		fmt.Fprintf(os.Stderr, "Error: Slot must be a number between 1 and 5\n")
		os.Exit(1)
	}

	// Find the project
	homeDir, _ := os.UserHomeDir()
	projectsDir := filepath.Join(homeDir, "projects")
	archiveDir := filepath.Join(homeDir, "archive")
	scratchDir := filepath.Join(homeDir, "scratch")

	projects, err := cache.FindProjectsCached(projectsDir, archiveDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to find projects: %v\n", err)
		os.Exit(1)
	}

	// Check scratch projects too
	scratchProjects, _ := findScratchProjects(scratchDir)
	projects = append(projects, scratchProjects...)

	// Find matching project
	var foundProject *config.Project
	for _, p := range projects {
		if strings.ToLower(p.ProjectInfo.ID) == projectName ||
			strings.ToLower(p.ProjectInfo.Name) == projectName {
			foundProject = p
			break
		}
	}

	if foundProject == nil {
		fmt.Fprintf(os.Stderr, "Error: Project '%s' not found\n", projectName)
		os.Exit(1)
	}

	// Check if slot is already occupied
	existingPin, _ := cache.GetPin(slot)
	if existingPin != nil {
		fmt.Printf("Replacing pin in slot %d: %s -> %s\n",
			slot, existingPin.ProjectID, foundProject.ProjectInfo.ID)
	}

	// Add the pin
	if err := cache.AddPin(slot, foundProject.ProjectInfo.ID, foundProject.Path); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to add pin: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("✓ Pinned '%s' to slot %d\n", foundProject.ProjectInfo.ID, slot)
	fmt.Printf("\nJump to it with:\n")
	fmt.Printf("  pk jump %d\n", slot)
	fmt.Printf("  Ctrl+b g %d\n", slot)
}

func runPinRemove(cmd *cobra.Command, args []string) {
	target := args[0]

	// Try to parse as slot number
	if slot, err := strconv.Atoi(target); err == nil {
		// Remove by slot
		pin, err := cache.GetPin(slot)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		if err := cache.RemovePin(slot); err != nil {
			fmt.Fprintf(os.Stderr, "Error: Failed to remove pin: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Removed pin from slot %d (%s)\n", slot, pin.ProjectID)
	} else {
		// Remove by project name
		if err := cache.RemovePinByProject(target); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Removed pin for project '%s'\n", target)
	}
}

func runPinList(cmd *cobra.Command, args []string) {
	pins, err := cache.ListPins()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to load pins: %v\n", err)
		os.Exit(1)
	}

	if len(pins) == 0 {
		fmt.Println("No pinned projects")
		fmt.Println("\nPin a project with:")
		fmt.Println("  pk pin add <project> <slot>")
		return
	}

	fmt.Println("Pinned Projects:")
	fmt.Println()

	for _, pin := range pins {
		fmt.Printf("  [%d]  %-20s  %s\n", pin.Slot, pin.ProjectID, pin.ProjectPath)
	}

	fmt.Println()
	fmt.Println("Jump to pinned projects with:")
	fmt.Println("  pk jump <slot>")
	fmt.Println("  Ctrl+b g <slot>")
}

func runPinClear(cmd *cobra.Command, args []string) {
	pins, _ := cache.ListPins()
	if len(pins) == 0 {
		fmt.Println("No pins to clear")
		return
	}

	fmt.Printf("This will remove all %d pin(s). Continue? (y/N): ", len(pins))
	var response string
	fmt.Scanln(&response)

	if strings.ToLower(response) != "y" {
		fmt.Println("Cancelled")
		return
	}

	if err := cache.ClearPins(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: Failed to clear pins: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("✓ All pins cleared")
}

// validPinAddArgs provides shell completion for pin add command
func validPinAddArgs(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	if len(args) == 0 {
		// First argument: project name
		return validAllProjectNames(cmd, args, toComplete)
	} else if len(args) == 1 {
		// Second argument: slot number
		return []string{"1", "2", "3", "4", "5"}, cobra.ShellCompDirectiveNoFileComp
	}
	return nil, cobra.ShellCompDirectiveNoFileComp
}
