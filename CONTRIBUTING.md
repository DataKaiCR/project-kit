# Contributing to PK

Thank you for your interest in contributing to PK (Project Kit).

## Getting Started

### Prerequisites

- Go 1.21 or later
- Optional: tmux, fzf (for session management features)
- Optional: cloud CLIs (for context switching features)

### Building from Source

```bash
git clone https://github.com/DataKaiCR/project-kit.git
cd project-kit
./build.sh
```

The binary will be in `bin/pk`.

## Development Workflow

### Making Changes

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/your-feature`)
3. Make your changes
4. Test your changes locally
5. Commit with clear, descriptive messages
6. Push to your fork
7. Open a pull request

### Commit Message Style

Follow the Linux kernel style:

- Short summary line (50 characters or less)
- Blank line
- Detailed explanation wrapped at 72 characters
- Focus on what and why, not how

Example:
```
Add support for git identity switching

Allows users to configure per-project git identities in .project.toml
using the git_identity field. Useful for maintaining separate work and
personal identities across projects.
```

### Code Style

- Follow standard Go conventions
- Run `go fmt` before committing
- Keep functions focused and modular
- Add comments for non-obvious logic
- Maintain the modular architecture (core + optional features)

## Project Structure

```
pk/
├── cmd/              # Command implementations (one file per command)
├── pkg/
│   ├── config/       # Configuration and .project.toml handling
│   ├── session/      # Tmux integration
│   ├── context/      # Cloud context switching
│   ├── cache/        # Project caching system
│   └── shell/        # Shell detection and alias generation
├── docs/             # Documentation (man pages, examples)
└── scripts/          # Installation and setup scripts
```

## Testing

While we don't currently have a comprehensive test suite, please ensure:

- Your changes work with the core commands
- Optional features gracefully degrade when dependencies are missing
- Changes don't break existing workflows

## Adding New Features

### Core Principles

1. **Modularity**: Core features should work standalone
2. **Graceful degradation**: Optional features should fail gracefully
3. **Minimal dependencies**: Don't add dependencies for core functionality
4. **User experience**: Prioritize clarity and ease of use

### Feature Guidelines

- New core commands should not require external dependencies
- Optional features should check for dependencies and provide helpful error messages
- Maintain backwards compatibility with existing .project.toml files
- Update documentation (README.md, man page) for new features

## Documentation

When adding features, update:

- README.md - Quick start and examples
- docs/pk.1 - Man page (roff format)
- Command help text (cobra short/long descriptions)
- Shell completion (if adding new commands)

## Questions or Issues

- Open an issue for bugs or feature requests
- Start a discussion for architectural questions
- Check existing issues before opening new ones

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
