package session

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/datakaicr/pk/pkg/config"
)

// CheckTmux verifies if tmux is installed
func CheckTmux() error {
	if _, err := exec.LookPath("tmux"); err != nil {
		return fmt.Errorf("'pk session' requires tmux to be installed\n" +
			"Install: brew install tmux (macOS) or apt install tmux (Linux)")
	}
	return nil
}

// IsInTmux checks if currently inside a tmux session
func IsInTmux() bool {
	return os.Getenv("TMUX") != ""
}

// SessionExists checks if a tmux session exists
func SessionExists(name string) bool {
	cmd := exec.Command("tmux", "has-session", "-t="+name)
	return cmd.Run() == nil
}

// SanitizeSessionName converts a project name to a valid tmux session name
func SanitizeSessionName(name string) string {
	// Replace dots with underscores (tmux doesn't like dots)
	return strings.ReplaceAll(name, ".", "_")
}

// CreateSession creates a new tmux session
func CreateSession(project *config.Project) error {
	sessionName := SanitizeSessionName(project.ProjectInfo.ID)

	// Check if session already exists
	if SessionExists(sessionName) {
		return SwitchSession(sessionName)
	}

	// Create new session based on configuration
	if len(project.Tmux.Windows) > 0 {
		return CreateWithLayout(project)
	}

	// Create basic session
	return CreateBasicSession(sessionName, project.Path)
}

// CreateBasicSession creates a simple single-window session
func CreateBasicSession(sessionName, path string) error {
	var cmd *exec.Cmd

	if IsInTmux() {
		// Inside tmux: create detached and switch
		cmd = exec.Command("tmux", "new-session", "-ds", sessionName, "-c", path)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to create tmux session: %w", err)
		}
		return SwitchSession(sessionName)
	}

	// Outside tmux: attach directly
	cmd = exec.Command("tmux", "new-session", "-s", sessionName, "-c", path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SwitchSession switches to an existing session
func SwitchSession(sessionName string) error {
	var cmd *exec.Cmd

	if IsInTmux() {
		cmd = exec.Command("tmux", "switch-client", "-t", sessionName)
	} else {
		cmd = exec.Command("tmux", "attach-session", "-t", sessionName)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	return cmd.Run()
}

// CreateWithLayout creates a session with custom window layout
func CreateWithLayout(project *config.Project) error {
	sessionName := SanitizeSessionName(project.ProjectInfo.ID)

	// Create base session (detached)
	cmd := exec.Command("tmux", "new-session", "-ds", sessionName, "-c", project.Path)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create session: %w", err)
	}

	// Kill the default window
	exec.Command("tmux", "kill-window", "-t", sessionName+":1").Run()

	// Create windows from configuration
	for i, window := range project.Tmux.Windows {
		windowPath := project.Path
		if window.Path != "" {
			windowPath = window.Path
		}

		windowName := window.Name
		if windowName == "" {
			windowName = fmt.Sprintf("window-%d", i+1)
		}

		// Create window
		windowTarget := fmt.Sprintf("%s:%d", sessionName, i+1)
		createCmd := exec.Command("tmux", "new-window", "-t", windowTarget, "-n", windowName, "-c", windowPath)
		if err := createCmd.Run(); err != nil {
			return fmt.Errorf("failed to create window %s: %w", windowName, err)
		}

		// Send command if specified
		if window.Command != "" {
			sendCmd := exec.Command("tmux", "send-keys", "-t", windowTarget, window.Command, "Enter")
			sendCmd.Run()
		}
	}

	// Set layout if specified
	if project.Tmux.Layout != "" {
		layoutCmd := exec.Command("tmux", "select-layout", "-t", sessionName, project.Tmux.Layout)
		layoutCmd.Run()
	}

	// Switch to session
	return SwitchSession(sessionName)
}

// ListSessions returns all active tmux sessions
func ListSessions() ([]string, error) {
	cmd := exec.Command("tmux", "list-sessions", "-F", "#{session_name}")
	output, err := cmd.Output()
	if err != nil {
		// No sessions is not an error
		return []string{}, nil
	}

	sessions := strings.Split(strings.TrimSpace(string(output)), "\n")
	return sessions, nil
}

// KillSession kills a tmux session by name
func KillSession(name string) error {
	cmd := exec.Command("tmux", "kill-session", "-t", name)
	return cmd.Run()
}
