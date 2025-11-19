package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var manCmd = &cobra.Command{
	Use:   "man",
	Short: "Open pk manual page",
	Long: `Open the pk man page in your system's man viewer.

If the man page is not installed system-wide, this will show
a helpful message with installation instructions.

Example:
  pk man`,
	Run: runMan,
}

func init() {
	rootCmd.AddCommand(manCmd)
}

func runMan(cmd *cobra.Command, args []string) {
	// Try to open system man page
	manCmd := exec.Command("man", "pk")
	manCmd.Stdin = os.Stdin
	manCmd.Stdout = os.Stdout
	manCmd.Stderr = os.Stderr

	if err := manCmd.Run(); err != nil {
		// Man page not installed, show helpful message
		fmt.Fprintf(os.Stderr, "Man page not installed.\n\n")
		fmt.Fprintf(os.Stderr, "To install:\n")
		fmt.Fprintf(os.Stderr, "  ./scripts/install-man.sh\n\n")
		fmt.Fprintf(os.Stderr, "Or view online:\n")
		fmt.Fprintf(os.Stderr, "  man docs/pk.1\n")
		os.Exit(1)
	}
}
