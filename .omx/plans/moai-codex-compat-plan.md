# MoAI-ADK for Codex: Reimplementation / Compatibility Layer Plan

## Requirements Summary

Goal: make `modu-ai/moai-adk` feel functionally complete inside Codex even though Claude Code specific internals cannot be ported 1:1.

Working definition of success:

- Users can keep the MoAI mental model: `init`, `project`, `plan`, `run`, `sync`, `review`, `coverage`, `clean`, `fix`, `loop`, `codemaps`, `worktree`, `status`, `doctor`, `update`.
- Project-local state and config remain MoAI-shaped under `.moai/` so existing docs and habits still apply.
- Claude-only surfaces are replaced with Codex-native equivalents, not dropped silently.
- In places where exact parity is impossible, the compatibility layer preserves outcome parity and command parity.

Key upstream facts this plan is based on:

- MoAI is a harness product, not just a CLI wrapper; its public value is the workflow surface around `plan -> run -> sync`, quality gates, session persistence, codemaps, clean/fix/loop, and agent orchestration (`README.ko.md`, sections around harness architecture and command list).
- The Go binary exposes a Cobra CLI rooted at `moai` and wires project/tool commands in `internal/cli/*.go`.
- `moai init` is a template/config/bootstrap pipeline, not a single config writer.
- Most user-visible behavior is generated from embedded templates under `internal/template/templates/`, especially `.claude/commands`, `.claude/skills`, `.claude/hooks`, `CLAUDE.md`, `.mcp.json`, and `.moai/config/sections/*`.
- Claude Code hook integration is deeply product-specific, so this is the main non-portable area.

## Recommended Strategy

Use a **Codex-native runtime with MoAI compatibility contracts**.

Decision:

- Keep the **outer contract** MoAI-compatible:
  - `moai` CLI entrypoint
  - `.moai/` directory layout
  - SPEC-oriented workflow names
  - status/doctor/update/worktree style commands
- Replace the **inner engine** with Codex/OMX-native execution:
  - `AGENTS.md` instead of `CLAUDE.md` as the orchestration brain
  - local Codex skills/prompts instead of `.claude/commands` and `.claude/skills`
  - OMX state/memory/team primitives instead of Claude hook events and Claude subagent APIs

Why this option:

- Upstream MoAI’s highest-value surface is in prompts/templates/config, not in the Go binary alone.
- A direct Claude-to-Codex port would spend most effort fighting product API differences in hooks and agent APIs.
- Preserving `.moai/` and command names keeps user experience stable while letting the execution layer be honest and maintainable.

Rejected options:

- Thin wrapper over upstream MoAI templates only: too brittle because generated assets target Claude Code semantics directly.
- Full byte-for-byte clone of Claude behavior: unrealistic and high-maintenance for low user value.
- Codex-only rewrite with no MoAI contract: easier technically, but loses the whole reason to emulate MoAI.

## Target Architecture

### 1. Public Compatibility Surface

Deliver these first-class user entrypoints:

- `moai init`
- `moai update`
- `moai doctor`
- `moai status`
- `moai project`
- `moai plan`
- `moai run`
- `moai sync`
- `moai review`
- `moai coverage`
- `moai clean`
- `moai fix`
- `moai loop`
- `moai codemaps`
- `moai worktree *`
- `moai cc`
- `moai cg`
- `moai glm`

Compatibility rule:

- If upstream behavior is Claude-specific, keep the command name and rebind it to the closest Codex/OMX outcome.

Examples:

- `moai cc`: Codex high-quality profile / leader-only mode.
- `moai glm`: cost-oriented alternate worker stack, likely Gemini-first where available.
- `moai cg`: Codex leader + cheaper parallel workers, preserving the “quality + cost balance” intent.

### 2. State and Config Layer

Preserve `.moai/` as the project contract:

- `.moai/config/sections/*`
- `.moai/project/*`
- `.moai/state/*`
- `.moai/manifest.json`

Add Codex-native internal assets alongside it:

- `AGENTS.md`
- `.codex/skills/moai*/` or equivalent project-local skill surface
- `.omx/` state for Codex execution bookkeeping

Rule:

- `.moai/` is the user-facing compatibility schema.
- `.omx/` is the internal execution schema.
- Generate or sync one from the other where needed; do not force users to learn a second model first.

### 3. Command Translation Layer

Build a `moai` router that maps MoAI commands onto Codex skills/modes:

- `plan` -> Codex `plan` skill with MoAI-flavored SPEC output
- `run` -> Codex execution workflow (`ralph` or `team`) with MoAI DDD/TDD guards
- `sync` -> docs/PR sync flow adapted to available git/GitHub tooling
- `review` -> Codex review workflow with MoAI quality checklist additions
- `fix` / `loop` -> Codex autonomous fix loop with explicit completion markers
- `project` / `codemaps` -> repo analysis and generated docs/flow artifacts

### 4. Hook Replacement Layer

Do **not** attempt a literal port of Claude hook registration.

Replace it with:

- command-time guards
- persistent state checks before/after workflow phases
- optional wrapper scripts for shell/git/worktree events
- OMX memory/state updates
- team-mode coordination hooks where OMX actually supports them

Principle:

- preserve the effect of hooks, not the transport mechanism.

## Planned Repository Structure

Recommended initial structure for this repo:

- `cmd/moai/main.go`
- `internal/cli/*`
- `internal/template/*`
- `internal/compat/*`
- `internal/workflows/*`
- `internal/spec/*`
- `internal/state/*`
- `templates/AGENTS.md.tmpl`
- `templates/.moai/config/sections/*`
- `templates/.codex/skills/moai/*`
- `templates/bin/*`
- `testdata/fixtures/*`

Design note:

- Reuse upstream package naming where it helps migration (`cli`, `template`, `manifest`, `worktree`, `status`).
- Introduce `compat` and `workflows` packages to isolate Codex translation from generic project bootstrap logic.

## Implementation Phases

### Phase 0. Contract Extraction

Purpose: freeze the user-visible MoAI contract before coding.

Tasks:

- Build a parity matrix from upstream commands, generated files, config sections, and workflow docs.
- Classify each item as:
  - native parity
  - adapted parity
  - stub/no-op with explanation
  - deferred
- Define explicit mappings for `cc`, `cg`, `glm`, hooks, and Claude-only skill routing.

Outputs:

- `docs/parity-matrix.md`
- `docs/command-mapping.md`
- `docs/non-portable-surfaces.md`

### Phase 1. Bootstrap CLI and Template Engine

Purpose: create the install/init/update backbone.

Tasks:

- Scaffold `moai` Go CLI with Cobra-compatible subcommands matching upstream names.
- Implement template deployment, manifest tracking, and safe re-update semantics.
- Preserve `.moai/manifest.json` and file provenance behavior.
- Generate `AGENTS.md`, `.moai/config/sections/*`, and project-local Codex skills in one `init`.

Acceptance criteria:

- `moai init .` creates a usable Codex project scaffold in an empty repo.
- `moai update` can re-sync templates without clobbering user-modified files.
- `moai doctor` validates Codex/OMX prerequisites and `.moai` health.

### Phase 2. Workflow Compatibility Commands

Purpose: make the high-value workflows usable end to end.

Tasks:

- Implement `project`, `plan`, `run`, `sync`, `review`, `coverage`, `clean`, `fix`, `loop`, `codemaps`.
- Generate MoAI-style SPECs and progress artifacts under `.moai/project/` and `.moai/state/`.
- Wrap Codex native skills (`plan`, `team`, `ralph`, `code-review`, `flowbook`, etc.) behind MoAI command vocabulary.

Acceptance criteria:

- A user can go from `moai init` to `moai plan` to `moai run` to `moai sync` without needing Claude-specific files.
- `plan` produces SPEC-like artifacts with testable acceptance criteria.
- `run` respects DDD/TDD mode from `.moai/config/sections/quality.yaml`.

### Phase 3. Team/Mode Emulation

Purpose: preserve MoAI’s multi-agent feel.

Tasks:

- Map MoAI execution modes to OMX capabilities:
  - solo -> single-agent Codex execution
  - team -> `team` / `ultrawork`
  - `cc` / `cg` / `glm` -> runtime presets
- Add mode selection rules mirroring upstream complexity thresholds where practical.
- Persist current mode in config/state for status output and resumability.

Acceptance criteria:

- `moai run --team` materially changes execution strategy.
- `moai cc|cg|glm` produce observable mode differences in status/config/output.

### Phase 4. Status, Worktree, and Recovery

Purpose: make the harness operable over time, not just on first run.

Tasks:

- Recreate `status`, `worktree`, `doctor`, `update`, and session/progress bookkeeping.
- Maintain resumable workflow state under `.moai/state/`.
- Add git-aware worktree helpers and safe branch/project switching affordances.

Acceptance criteria:

- `moai status` shows branch, mode, active workflow, and health summary.
- interrupted `run` / `loop` workflows can resume from persisted state.
- worktree commands operate safely in normal git repos.

### Phase 5. Gap Closure and Experience Polish

Purpose: remove remaining “this is obviously not MoAI” edges.

Tasks:

- Improve prompt wording and generated docs so they read like MoAI, not generic Codex scaffolding.
- Add migration helpers for existing `.moai/` projects.
- Add explicit user-facing notices for adapted or degraded features.
- Validate on empty repo, existing app repo, and multi-language repo scenarios.

Acceptance criteria:

- A MoAI user can use the scaffold without reading internal compatibility docs first.
- All known non-parity cases are surfaced intentionally, not as surprise failures.

## Command Mapping Plan

### High-priority exact-or-close parity

Implement first:

- `init`
- `project`
- `plan`
- `run`
- `sync`
- `review`
- `coverage`
- `clean`
- `fix`
- `loop`
- `codemaps`
- `status`
- `doctor`
- `update`
- `worktree`

### Adapted parity

Implement with explicit remapping:

- `cc`
- `cg`
- `glm`
- Claude hook commands
- Claude-specific `.claude/*` generated assets

### Likely defer or narrow

- OAuth/MCP behaviors that only exist because Claude Code consumes `.mcp.json` directly
- provider-specific assumptions tied to Anthropic/GLM launch semantics

## File-Level Plan

Initial files to create here:

- `/home/pi/code/moai-adk-codex/cmd/moai/main.go`
- `/home/pi/code/moai-adk-codex/internal/cli/root.go`
- `/home/pi/code/moai-adk-codex/internal/cli/init.go`
- `/home/pi/code/moai-adk-codex/internal/cli/update.go`
- `/home/pi/code/moai-adk-codex/internal/cli/status.go`
- `/home/pi/code/moai-adk-codex/internal/cli/doctor.go`
- `/home/pi/code/moai-adk-codex/internal/cli/worktree/*.go`
- `/home/pi/code/moai-adk-codex/internal/template/deployer.go`
- `/home/pi/code/moai-adk-codex/internal/template/embed.go`
- `/home/pi/code/moai-adk-codex/internal/compat/router.go`
- `/home/pi/code/moai-adk-codex/internal/compat/modes.go`
- `/home/pi/code/moai-adk-codex/internal/workflows/*.go`
- `/home/pi/code/moai-adk-codex/templates/AGENTS.md.tmpl`
- `/home/pi/code/moai-adk-codex/templates/.moai/config/sections/*.tmpl`
- `/home/pi/code/moai-adk-codex/templates/.codex/skills/moai/SKILL.md`

## Acceptance Criteria

- A fresh repo can run `moai init .` and receive a Codex-usable scaffold plus MoAI-compatible `.moai/` state.
- Generated assets contain no mandatory Claude-only dependency on `.claude/commands`, Claude hook transport, or Claude-specific agent APIs.
- `moai plan`, `moai run`, and `moai sync` each perform useful Codex-native work, not placeholder output.
- The compatibility layer documents every intentionally adapted behavior.
- User-modified managed files survive `moai update`.
- Team mode, solo mode, and at least one cost-optimized mode are observable and test-covered.

## Risks and Mitigations

### Risk: copying names without preserving outcomes

Mitigation:

- build parity matrix first
- require command-level behavior tests
- define “adapted parity” explicitly

### Risk: Claude hook assumptions leaking into Codex implementation

Mitigation:

- isolate hook replacement logic in `internal/compat/*`
- ban direct generation of unusable `.claude/hooks/*` as required runtime dependencies

### Risk: `.moai` and `.omx` state drifting

Mitigation:

- define a single source of truth per concern
- add sync tests for mode, progress, and manifest state

### Risk: `cc` / `cg` / `glm` semantics becoming misleading

Mitigation:

- rename internally, preserve external names only
- print active mapped mode in status output
- document the mapping in generated project docs

### Risk: over-building low-value parity

Mitigation:

- deliver in the priority order above
- defer low-value Claude-only edge cases until core workflow parity is stable

## Verification Steps

- Unit tests for CLI routing, manifest tracking, template deployment, mode mapping, and state sync.
- Golden tests for generated templates from `init`.
- Integration tests:
  - empty repo -> `init`
  - initialized repo -> `update`
  - sample repo -> `project`, `plan`
  - SPEC workflow -> `run`, `sync`
- UX verification:
  - confirm a MoAI user can discover the same top-level commands
  - confirm status/doctor explain adapted surfaces clearly

## Upstream Evidence References

- Command root and CLI framing:
  - `/tmp/moai-adk/internal/cli/root.go`
  - `/tmp/moai-adk/cmd/moai/main.go`
- Init/bootstrap pipeline:
  - `/tmp/moai-adk/internal/cli/init.go`
- Hook/event transport:
  - `/tmp/moai-adk/internal/cli/hook.go`
  - `/tmp/moai-adk/internal/template/templates/.claude/settings.json.tmpl`
- Template/manifest deployment contract:
  - `/tmp/moai-adk/internal/template/deployer.go`
  - `/tmp/moai-adk/internal/template/templates/`
- User-facing workflow and command surface:
  - `/tmp/moai-adk/README.ko.md`

## Recommendation for Execution Order

1. Phase 0 and Phase 1 first, in the same PR or milestone.
2. Then Phase 2 for `project/plan/run/sync`.
3. Then Phase 3 and Phase 4 for team modes and operability.
4. Leave Phase 5 as targeted polish driven by parity gaps found in real use.
