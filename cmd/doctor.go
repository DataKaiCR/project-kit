package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/datakaicr/pk/pkg/cache"
	"github.com/datakaicr/pk/pkg/paths"
	"github.com/spf13/cobra"
)

var doctorCmd = &cobra.Command{
	Use:   "doctor",
	Short: "Diagnose pk setup and configuration issues",
	Long: `Check pk installation, dependencies, and configuration for common issues.

This command performs health checks on:
  - Directory structure (~/projects, ~/archive, ~/scratch)
  - Dependencies (tmux, fzf)
  - Tmux configuration
  - Cache file integrity
  - Stale path detection
  - Config file validity

Example:
  pk doctor`,
	Run: runDoctor,
}

func init() {
	rootCmd.AddCommand(doctorCmd)
}

func runDoctor(cmd *cobra.Command, args []string) {
	fmt.Println("PK Doctor - Diagnosing your setup...")
	fmt.Println()

	issues := 0

	// Check 1: Path resolver and directory structure
	fmt.Println("📂 Checking directory structure...")
	resolver, err := paths.NewResolver()
	if err != nil {
		fmt.Printf("   ❌ Failed to create path resolver: %v\n", err)
		issues++
	} else {
		checkDirectory(resolver.Projects(), "Projects directory", &issues)
		checkDirectory(resolver.Archive(), "Archive directory", &issues)
		checkDirectory(resolver.Scratch(), "Scratch directory", &issues)
		checkDirectory(resolver.Scriptorium(), "Scriptorium directory", &issues)
	}
	fmt.Println()

	// Check 2: Dependencies
	fmt.Println("🔧 Checking dependencies...")
	checkCommand("tmux", "Required for 'pk session' and tmux keybindings", &issues)
	checkCommand("fzf", "Required for interactive project selection", &issues)
	fmt.Println()

	// Check 3: Tmux configuration
	fmt.Println("⚙️  Checking tmux configuration...")
	checkTmuxConfig(&issues)
	fmt.Println()

	// Check 4: Cache integrity
	fmt.Println("💾 Checking cache files...")
	checkCacheIntegrity(&issues)
	fmt.Println()

	// Check 5: Config file
	fmt.Println("📝 Checking configuration...")
	checkConfigFile(&issues)
	fmt.Println()

	// Check 6: Stale paths
	fmt.Println("🔍 Checking for stale paths...")
	checkStalePaths(&issues)
	fmt.Println()

	// Summary
	fmt.Println("════════════════════════════════════════")
	if issues == 0 {
		fmt.Println("✅ All checks passed! PK is healthy.")
	} else {
		fmt.Printf("⚠️  Found %d issue(s) that need attention.\n", issues)
	}
	fmt.Println("════════════════════════════════════════")
}

func checkDirectory(path, name string, issues *int) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("   ⚠️  %s does not exist: %s\n", name, path)
		fmt.Printf("      Run: mkdir -p %s\n", path)
		*issues++
	} else {
		fmt.Printf("   ✓ %s: %s\n", name, path)
	}
}

func checkCommand(name, description string, issues *int) {
	if _, err := exec.LookPath(name); err == nil {
		fmt.Printf("   ✓ %s installed\n", name)
	} else {
		fmt.Printf("   ❌ %s not found - %s\n", name, description)
		fmt.Printf("      Install: apt install %s (Debian/Ubuntu) or brew install %s (macOS)\n", name, name)
		*issues++
	}
}

func checkTmuxConfig(issues *int) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("   ❌ Cannot determine home directory\n")
		*issues++
		return
	}

	// Check both possible tmux config locations
	configPaths := []string{
		filepath.Join(homeDir, ".config", "tmux", "tmux.conf"),
		filepath.Join(homeDir, ".tmux.conf"),
	}

	configFound := false
	var foundPath string

	for _, path := range configPaths {
		if _, err := os.Stat(path); err == nil {
			configFound = true
			foundPath = path
			break
		}
	}

	if !configFound {
		fmt.Printf("   ⚠️  No tmux config found\n")
		fmt.Printf("      Tmux keybindings (Ctrl+b f, Ctrl+b g) will not work\n")
		fmt.Printf("      See: docs/tmux-keybindings.conf\n")
		*issues++
		return
	}

	// Check if PK keybindings are in config
	data, err := os.ReadFile(foundPath)
	if err != nil {
		fmt.Printf("   ⚠️  Cannot read tmux config: %v\n", err)
		*issues++
		return
	}

	configContent := string(data)
	hasPkSession := containsString(configContent, "pk session") || containsString(configContent, "pk sessions")
	hasPkJump := containsString(configContent, "pk jump")

	if hasPkSession && hasPkJump {
		fmt.Printf("   ✓ Tmux config found with PK keybindings: %s\n", foundPath)
	} else {
		fmt.Printf("   ⚠️  Tmux config exists but missing PK keybindings: %s\n", foundPath)
		fmt.Printf("      Add keybindings from: docs/tmux-keybindings.conf\n")
		*issues++
	}
}

func checkCacheIntegrity(issues *int) {
	cacheFile, err := cache.GetCacheFile()
	if err != nil {
		fmt.Printf("   ❌ Cannot determine cache location: %v\n", err)
		*issues++
		return
	}

	if _, err := os.Stat(cacheFile); os.IsNotExist(err) {
		fmt.Printf("   ℹ️  Cache not yet built (will be created on first use)\n")
	} else {
		fmt.Printf("   ✓ Cache file exists: %s\n", cacheFile)

		// Try to load cache
		if _, err := cache.LoadFromCache(); err != nil {
			fmt.Printf("   ❌ Cache file corrupted: %v\n", err)
			fmt.Printf("      Run: pk cache clear && pk cache refresh\n")
			*issues++
		}
	}
}

func checkConfigFile(issues *int) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Printf("   ❌ Cannot determine home directory\n")
		*issues++
		return
	}

	configPath := filepath.Join(homeDir, ".config", "pk", "config.toml")

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		fmt.Printf("   ℹ️  No config file (using defaults)\n")
		fmt.Printf("      Optional: Create %s to customize paths\n", configPath)
	} else {
		// Try to load config
		_, err := paths.NewResolver()
		if err != nil {
			fmt.Printf("   ❌ Config file exists but has errors: %s\n", configPath)
			fmt.Printf("      Error: %v\n", err)
			*issues++
		} else {
			fmt.Printf("   ✓ Config file loaded: %s\n", configPath)
		}
	}
}

func checkStalePaths(issues *int) {
	resolver, err := paths.NewResolver()
	if err != nil {
		fmt.Printf("   ❌ Cannot check paths: %v\n", err)
		*issues++
		return
	}

	staleCount := 0

	// Check pins
	pins, err := cache.LoadPins()
	if err == nil {
		for _, pin := range pins {
			if _, err := os.Stat(pin.ProjectPath); os.IsNotExist(err) {
				// Try to find it
				if _, err := resolver.FindProject(pin.ProjectID); err != nil {
					fmt.Printf("   ⚠️  Pin [%d] %s: path not found\n", pin.Slot, pin.ProjectID)
					staleCount++
				}
			}
		}
	}

	// Check access records
	records, err := cache.LoadAccessRecords()
	if err == nil {
		for _, record := range records {
			if _, err := os.Stat(record.ProjectPath); os.IsNotExist(err) {
				// Try to find it
				if _, err := resolver.FindProject(record.ProjectID); err != nil {
					staleCount++
				}
			}
		}
	}

	if staleCount == 0 {
		fmt.Printf("   ✓ All cached paths are valid\n")
	} else {
		fmt.Printf("   ℹ️  Found %d stale path(s) - they will be auto-healed on next use\n", staleCount)
	}
}

func containsString(haystack, needle string) bool {
	return len(haystack) >= len(needle) &&
		   (haystack == needle ||
		    haystack[:len(needle)] == needle ||
		    containsString(haystack[1:], needle))
}
