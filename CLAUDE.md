# PK - Project Kit Development Instructions

## Session Start

**CRITICAL**: On every session start, immediately call `dk_get_context_snapshot()` to load project context automatically.

This provides:
- Current project metadata
- Last commit and current branch
- Files changed since last commit
- Current focus area
- Applicable protocols
- Related documentation

Do NOT skip this step. DKOS is built on automated context management via MCP tools.

## Project Overview

Command-line tool for managing software projects with metadata tracking, tmux session management, and cloud context switching. Features include project lifecycle management, scratch projects, git repository cloning, recent project tracking, and shell alias generation. Built with Go for cross-platform support.

**Project Type:** internal
**Status:** active
**Tech Stack:** go, cli

## Architecture

This project is part of the DataKai ecosystem:
- Follows DataKai development protocols
- Uses DKOS MCP tools for development intelligence
- Integrates with Chronicle for cross-project documentation

See README.md and docs/architecture/ for details.

## Development Protocols

All commits are validated via dkproto binary with context-aware variants:
- DataKai internal projects: `type(scope): description` format
- Enforces: No DRY runs, Semantic Self-Containment, proper documentation

## Available MCP Tools

Use these tools proactively:
- `dk_get_context_snapshot()` - Load project context
- `dk_query_docs(query)` - Search DataKai documentation
- `dk_get_standards(category)` - Query protocol standards
- `dk_check_commit(message)` - Validate commit messages
- `dk_install_hooks()` - Install git hooks

## Notes

- Eye documentation lives in `~/projects/eye/` (strategic vision source of truth)
- dkproto binary lives in `~/projects/dkproto/` (protocol kernel)
- This project uses PK metadata (.project.toml) for context discovery

- Repository: https://github.com/datakaicr/pk
