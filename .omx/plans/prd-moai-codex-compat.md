# PRD: MoAI Codex Compatibility Layer

## Goal

Deliver a Codex-native implementation that preserves the user-facing MoAI workflow surface.

## User Stories

### US-001 Bootstrap a compatible project

As a MoAI user, I want `moai init` to scaffold a project for Codex so that I can start using MoAI-style workflows without Claude-specific runtime dependencies.

Acceptance criteria:

- `moai init .` creates `.moai/` config/state scaffolding.
- `moai init .` generates `AGENTS.md` and project-local command wrappers/templates.
- Existing files are not overwritten unsafely on re-run.

### US-002 Operate the compatibility layer

As a MoAI user, I want `status`, `doctor`, and `update` commands so that I can inspect and maintain the scaffold over time.

Acceptance criteria:

- `moai status` reports project, mode, and scaffold health.
- `moai doctor` reports missing prerequisites and scaffold issues.
- `moai update` safely reapplies managed templates.

### US-003 Keep the workflow vocabulary

As a MoAI user, I want `project`, `plan`, `run`, `sync`, `review`, `coverage`, `clean`, `fix`, `loop`, and `codemaps` commands so that the familiar workflow vocabulary still exists inside Codex.

Acceptance criteria:

- Each command is routable through the CLI.
- Each command performs a meaningful Codex-compatible action or writes an explicit compatibility artifact.
- The compatibility behavior is documented.

### US-004 Ship with documentation

As a repository user, I want a clear README so that I can install, use, and understand the compatibility layer.

Acceptance criteria:

- `README.md` explains goals, constraints, install/use, and current parity.
- README matches actual implemented commands.
