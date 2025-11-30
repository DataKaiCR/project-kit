# Protocol Kernel Projects (PKP)

## Overview

A **Protocol Kernel Project (PKP)** is a "smart project"—a directory that has been elevated from passive storage to an active participant in an intelligent development ecosystem.

## The PKP Trinity

Every PKP has three layers:

```
┌─────────────────────────────────────────┐
│     Communication (Ubiweave)            │
│     Cross-project messaging             │
└─────────────────────────────────────────┘
                    ↑
┌─────────────────────────────────────────┐
│     Intelligence (DKOS Hooks)           │
│     SessionStart, PreCompact, etc.      │
└─────────────────────────────────────────┘
                    ↑
┌─────────────────────────────────────────┐
│     Identity (.project.toml)            │
│     Metadata, context, relationships    │
└─────────────────────────────────────────┘
```

### Layer 1: Identity

The `.project.toml` file is the birth certificate of a PKP. It provides:

- **Metadata**: Name, ID, status, type
- **Tech context**: Stack, domain, frameworks
- **Relationships**: Links to repos, docs, knowledge graphs
- **Configuration**: Tmux layouts, cloud contexts, git identities

Without identity, a project is just a directory. With identity, it becomes discoverable, queryable, and manageable.

**Created by**: `pk new`, `pk clone`, `pk promote`

### Layer 2: Intelligence

DKOS hooks inject intelligence into the development workflow:

| Hook | Trigger | Purpose |
|------|---------|---------|
| `SessionStart` | `pk session <project>` | Load context, check health, surface relevant info |
| `PreCompact` | Context approaching limit | Create CCC file, preserve continuity |
| `user_prompt_submit` | Every prompt | Protocol validation, context injection |
| `commit-msg` | Git commit | Validate against project protocols |

Intelligence transforms a project from passive storage to active assistant.

**Enabled by**: `dk_bootstrap`, DKOS MCP server

### Layer 3: Communication

Ubiweave enables PKPs to communicate across project boundaries:

- **Post bugs** from one project to another's inbox
- **Broadcast events** (releases, breaking changes)
- **Query knowledge** from related projects
- **Coordinate workflows** across the ecosystem

Communication transforms isolated projects into a collaborative mesh.

**Enabled by**: Ubiweave service mesh (Phase 1: file-based, Phase 2: event bus)

## PKP Maturity Levels

Not all projects need all three layers. PKP maturity is progressive:

| Level | Layers | Description |
|-------|--------|-------------|
| **L0** | None | Plain directory |
| **L1** | Identity | Has `.project.toml`, managed by pk |
| **L2** | Identity + Intelligence | DKOS hooks active, protocol validation |
| **L3** | Full PKP | All three layers, Ubiweave connected |

Most open source projects stay at L1. DataKai internal projects operate at L2-L3.

## The pk Command: Dual Identity

The `pk` command serves two audiences:

### For Open Source Users: Project Kit

A CLI tool for project lifecycle management:
- Create, clone, archive projects
- Track metadata in `.project.toml`
- Manage tmux sessions
- Switch cloud contexts

No intelligence, no communication—just organized project management.

### For DataKai Ecosystem: Protocol Kernel

The foundation layer for PKPs:
- Creates identity (`.project.toml`)
- Integrates with DKOS for intelligence
- Connects to Ubiweave for communication
- Enforces protocol compliance

Same tool, deeper integration.

## Ubiweave: The Communication Layer

Ubiweave is the service mesh that enables cross-PKP communication.

### Channels

Messages are routed through semantic channels:

| Channel | Purpose | Example |
|---------|---------|---------|
| `infra.*` | Infrastructure events | `infra.deploy.production` |
| `knowledge.*` | Knowledge updates | `knowledge.zettel.created` |
| `strategic.*` | Strategic decisions | `strategic.architecture.changed` |
| `docs.*` | Documentation events | `docs.readme.updated` |
| `broadcast.*` | Ecosystem-wide announcements | `broadcast.release.dkos-1.0` |

### Phase 1: File-Based Queue

Initial implementation uses filesystem:

```
.context/
├── incoming/           # Messages received
│   ├── bug-001.json
│   └── task-002.json
├── outgoing/           # Messages to send
└── processed/          # Archived messages
```

Tools:
- `dk_post_to_project()` - Send message to another PKP
- `dk_check_inbox()` - Check for incoming messages
- `dk_promote_inbox()` - Convert message to roadmap task

### Phase 2: Event Bus (Future)

When file-based becomes limiting:
- Redis/NATS for real-time messaging
- Event sourcing for audit trail
- Subscription-based routing

## PKP in Practice

### Creating a PKP (L1)

```bash
pk new my-project
cd ~/projects/my-project
# .project.toml created with identity
```

### Elevating to L2

```bash
dk_bootstrap  # Install DKOS hooks
# Now has SessionStart, PreCompact, commit validation
```

### Full PKP (L3)

```bash
# Project participates in Ubiweave
dk_post_to_project \
  --target dkos \
  --type bug \
  --title "Session hook not firing" \
  --body "Details..."

# Check for incoming messages
dk_check_inbox
```

## Architecture Integration

PKPs exist within the broader DataKai architecture:

```
┌─────────────────────────────────────────────────────────┐
│                    Products (Layer 4)                    │
│            DataDojo, Pipelines, Navigator               │
└─────────────────────────────────────────────────────────┘
                           ↑
┌─────────────────────────────────────────────────────────┐
│              Knowledge Graphs (Layer 3)                  │
│         Data Engineering, GenAI, Cloud KGs              │
└─────────────────────────────────────────────────────────┘
                           ↑
┌─────────────────────────────────────────────────────────┐
│           Knowledge Infrastructure (Layer 2)             │
│                      Conduit                             │
└─────────────────────────────────────────────────────────┘
                           ↑
┌─────────────────────────────────────────────────────────┐
│           Intelligence Services (Layer 1)                │
│              DKOS, Scriptoria, Chronicle                 │
└─────────────────────────────────────────────────────────┘
                           ↑
┌─────────────────────────────────────────────────────────┐
│              Protocol Kernel (Layer 0)                   │
│                      dkproto                             │
└─────────────────────────────────────────────────────────┘
                           ↑
┌─────────────────────────────────────────────────────────┐
│           Protocol Kernel Projects (PKPs)                │
│    Identity (pk) + Intelligence (DKOS) + Comms (Ubiweave)│
└─────────────────────────────────────────────────────────┘
```

PKPs are the **ground level**—the individual projects that generate knowledge, receive intelligence, and communicate through the mesh.

## Why PKP Matters

### For Individual Developers

- Projects remember their context
- AI assistants understand project history
- Switching projects is seamless

### For Teams

- Shared project metadata
- Consistent protocol enforcement
- Cross-project visibility

### For the Ecosystem

- Projects communicate autonomously
- Knowledge flows between related work
- The system learns from all projects

## Comparison: Traditional vs PKP

| Aspect | Traditional Project | Protocol Kernel Project |
|--------|---------------------|------------------------|
| Identity | README maybe | `.project.toml` with rich metadata |
| Context | Lost between sessions | Preserved via CCC files |
| Protocols | Ad-hoc, inconsistent | Enforced by dkproto |
| AI Integration | Generic assistance | Context-aware intelligence |
| Cross-project | Manual coordination | Ubiweave messaging |
| Discovery | grep/find | `pk list`, `dk_query_docs` |

## Implementation Status

| Component | Status | Notes |
|-----------|--------|-------|
| Identity (pk) | ✅ Production | `.project.toml` schema complete |
| Intelligence (DKOS) | ✅ Active | Hooks implemented |
| Communication (Ubiweave) | 🔄 Phase 1 | File-based inbox working |
| Full PKP | 🎯 Target | All three layers integrated |

## Related Documentation

- [Schema Design](schema-design.md) - `.project.toml` structure
- [DKOS Architecture](../../dkos/README.md) - Intelligence layer
- [DataKai Architecture](../../eye/docs/architecture/OVERVIEW.md) - Ecosystem context

## Summary

A **Protocol Kernel Project** is:

1. **Identified** by `.project.toml` (created by pk)
2. **Intelligent** through DKOS hooks (session management, protocol validation)
3. **Connected** via Ubiweave (cross-project communication)

The `pk` command is the entry point—it creates project identity. DKOS adds intelligence. Ubiweave enables communication. Together, they transform directories into smart, connected participants in the DataKai ecosystem.
