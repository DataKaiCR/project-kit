# PK - Project Kit

Command-line tool for managing software projects with metadata tracking, tmux session management, and cloud context switching.

## Installation

### Quick Install

```bash
git clone https://github.com/datakaicr/project-kit.git
cd project-kit
make install
pk install
```

This creates the required directories (`~/projects`, `~/scratch`, `~/archive`) and installs the binary, man page, and shell completions.

### Manual Build

```bash
./build.sh
./bin/pk install
```

### Requirements

**Core:** Go 1.21+ (build only)

**Optional:**
- tmux, fzf - for session management
- aws, az, gcloud - for context switching

## Quick Start

```bash
# Create experimental project
pk scratch new api-test
cd ~/scratch/api-test

# Promote to real project when ready
pk promote api-test

# Open in tmux session
pk session api-test
```

## Core Commands

### Project Management

```bash
pk new <name>              # Create project in ~/projects
pk clone <url> [name]      # Clone git repo and create .project.toml
pk list [filter]           # List projects (active, archived, etc.)
pk show <name>             # View project details
pk recent                  # List recently accessed projects
pk edit <name>             # Edit metadata
pk rename <old> <new>      # Rename project
pk archive <name>          # Move to ~/archive
pk delete <name>           # Remove permanently

pk pin add <name> <slot>   # Pin project to slot (1-5)
pk pin list                # List pinned projects
pk jump <slot>             # Jump to pinned project
```

### Scratch Projects

Lightweight projects for experimentation in `~/scratch`.

```bash
pk scratch new <name>      # Create scratch project
pk scratch list            # View all scratch projects
pk scratch delete <name>   # Remove scratch project
pk promote <name>          # Convert to full project
```

### Session Management

Requires tmux and fzf.

```bash
pk session                 # Interactive project selector (all projects)
pk session <name>          # Open specific project
pk sessions                # Active sessions only (fast, Harpoon-style)
pk sessions <name>         # Switch to active session directly
```

Features:
- Custom window layouts
- Active session indicators
- tmux configuration via `.project.toml`

### Pinned Projects (Harpoon-style)

Pin your most-used projects to numbered slots for instant access:

```bash
pk pin add pk 1            # Pin 'pk' to slot 1
pk pin add dkos 2          # Pin 'dkos' to slot 2
pk pin add conduit 3       # Pin 'conduit' to slot 3
pk pin list                # Show all pins
pk pin remove 1            # Remove pin from slot 1

pk jump 1                  # Jump to slot 1 (opens tmux session)
pk jump 2                  # Jump to slot 2
```

**Tmux Keybindings:**

Add to `~/.tmux.conf` for instant jumping:

```bash
# Session jumping with Ctrl+b g <number>
bind-key g switch-client -T jump
bind-key -T jump 1 run-shell "pk jump 1"
bind-key -T jump 2 run-shell "pk jump 2"
bind-key -T jump 3 run-shell "pk jump 3"
bind-key -T jump 4 run-shell "pk jump 4"
bind-key -T jump 5 run-shell "pk jump 5"

# Fast active session switcher (Ctrl+b F)
bind-key F run-shell "tmux display-popup -E -w 90% -h 80% 'pk sessions'"
```

Then use:
- `Ctrl+b g 1` - Jump to pinned project in slot 1
- `Ctrl+b g 2` - Jump to pinned project in slot 2
- `Ctrl+b F` - Switch between active sessions only

### Aliases

```bash
pk sync                    # Generate shell aliases
```

Creates aliases like `dojo` to jump to projects. Run after creating or renaming projects.

## Project Metadata

Projects use `.project.toml` for metadata:

```toml
[project]
name = "My Project"
id = "my-project"
status = "active"
type = "product"

[ownership]
primary = "owner-name"

[dates]
started = "2025-01-15"
```

See `docs/examples/` for complete configuration examples.

### Tmux Configuration

```toml
[tmux]
layout = "main-vertical"
windows = [
    {name = "editor", command = "nvim"},
    {name = "server", command = "npm run dev"},
    {name = "logs", command = "tail -f logs/app.log"}
]
```

### Context Switching

```toml
[context]
aws_profile = "production"
azure_subscription = "My Subscription"
gcloud_project = "my-gcp-project"
databricks_profile = "prod"
git_identity = "work"
```

When opening a session, pk automatically switches to configured contexts.

## Architecture

```
pk (Core - No Dependencies)
├── Project lifecycle management
├── Metadata tracking
├── Caching
└── Shell alias generation

Optional Modules
├── Session (requires tmux, fzf)
│   ├── Project switching
│   └── Custom layouts
└── Context (requires cloud CLIs)
    ├── AWS, Azure, GCloud
    └── Git identity switching
```

Core commands work standalone. Optional features require external tools.

## Common Workflows

### Experimentation to Production

```bash
pk scratch new prototype
cd ~/scratch/prototype
# ... experiment ...
pk promote prototype
pk edit prototype
# ... add metadata ...
pk session prototype
```

### Project Navigation

```bash
pk list active             # View active projects
pk recent                  # View recently accessed projects
pk session                 # Interactive tmux selector
pk show myproject          # View details
```

### Cloning Projects

```bash
pk clone https://github.com/user/repo
pk clone git@github.com:user/repo.git my-name
pk clone https://github.com/user/repo --session  # Clone and open
```

### Shell Aliases

```bash
pk sync                    # Generate aliases
source ~/.zshrc            # Reload shell
myproject                  # Jump to project
```

## Integration

### Neovim

```lua
-- Add to init.lua
vim.keymap.set('n', '<C-f>', function()
    vim.fn.system('tmux display-popup -E "pk session"')
end)
```

### Shell Completion

Installed automatically with `pk install` or manually:

```bash
pk completion zsh > /usr/local/share/zsh/site-functions/_pk
pk completion bash > ~/.bash_completion.d/pk
pk completion fish > ~/.config/fish/completions/pk.fish
```

## File Locations

```
~/.cache/pk/projects.json              # Project cache (5min TTL)
~/.config/zsh/project-aliases.zsh      # Shell aliases (zsh)
~/.bash_aliases                        # Shell aliases (bash)
~/.config/fish/conf.d/project-aliases.fish  # Shell aliases (fish)
```

## Documentation

```bash
pk --help                  # Command help
pk <command> --help        # Command-specific help
man pk                     # Full manual
```

## Development

```bash
# Build
make build

# Install locally
make install

# Clean
make clean

# Run tests
make test
```

## Project Structure

```
pk/
├── cmd/              # Command implementations
├── pkg/
│   ├── config/       # .project.toml handling
│   ├── session/      # Tmux integration
│   ├── context/      # Cloud context switching
│   ├── cache/        # Project caching
│   └── shell/        # Alias generation
├── docs/
│   └── pk.1          # Man page
├── scripts/          # Install scripts
├── build.sh          # Build script
└── Makefile          # Build automation
```

## License

Open source. See LICENSE file.

## Contributing

Contributions welcome. Open an issue or pull request at:
https://github.com/datakaicr/pk
