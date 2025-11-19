package config

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
)

// Project represents a .project.toml file
type Project struct {
	Path string // Full path to project directory

	// [project] section
	ProjectInfo struct {
		Name   string `toml:"name"`
		ID     string `toml:"id"`
		Status string `toml:"status"`
		Type   string `toml:"type"`
	} `toml:"project"`

	// [ownership] section
	Ownership struct {
		Primary      string   `toml:"primary"`
		Partners     []string `toml:"partners"`
		LicenseModel string   `toml:"license_model"`
	} `toml:"ownership"`

	// [client] section
	Client struct {
		EndClient    string `toml:"end_client"`
		Intermediary string `toml:"intermediary"`
		MyRole       string `toml:"my_role"`
	} `toml:"client"`

	// [tech] section
	Tech struct {
		Stack  []string `toml:"stack"`
		Domain []string `toml:"domain"`
	} `toml:"tech"`

	// [dates] section
	Dates struct {
		Started   string `toml:"started"`
		Completed string `toml:"completed"`
	} `toml:"dates"`

	// [links] section
	Links struct {
		ScriptoriumProject string `toml:"scriptorium_project"`
		Repository         string `toml:"repository"`
		Documentation      string `toml:"documentation"`
		ConduitGraph       string `toml:"conduit_graph"`
	} `toml:"links"`

	// [notes] section
	Notes struct {
		Description string `toml:"description"`
	} `toml:"notes"`

	// [tmux] section (optional - for pk session)
	Tmux struct {
		Layout  string         `toml:"layout"`
		Windows []TmuxWindow   `toml:"windows"`
	} `toml:"tmux"`

	// [context] section (optional - for pk context)
	Context struct {
		AWSProfile       string `toml:"aws_profile"`
		AzureSubscription string `toml:"azure_subscription"`
		GCloudProject    string `toml:"gcloud_project"`
		DatabricksProfile string `toml:"databricks_profile"`
		SnowflakeAccount string `toml:"snowflake_account"`
		GitIdentity      string `toml:"git_identity"`
	} `toml:"context"`
}

// TmuxWindow represents a window configuration
type TmuxWindow struct {
	Name    string `toml:"name"`
	Command string `toml:"command"`
	Path    string `toml:"path"`
}

// LoadProject reads a .project.toml file
func LoadProject(path string) (*Project, error) {
	var project Project
	project.Path = filepath.Dir(path)

	// Decode TOML file
	if _, err := toml.DecodeFile(path, &project); err != nil {
		return nil, err
	}

	return &project, nil
}

// FindProjects recursively finds all .project.toml files
func FindProjects(rootDirs ...string) ([]*Project, error) {
	var projects []*Project

	for _, root := range rootDirs {
		// Check if directory exists
		if _, err := os.Stat(root); os.IsNotExist(err) {
			continue
		}

		// Walk directory tree
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			// Found a .project.toml file
			if info.Name() == ".project.toml" {
				project, err := LoadProject(path)
				if err != nil {
					// Skip malformed files
					return nil
				}
				projects = append(projects, project)
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return projects, nil
}
