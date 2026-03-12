# moai-adk-codex

MoAI-ADK compatibility layer for Codex and OMX, exposed as the `moai-codex` command.

This project preserves the familiar MoAI command surface while swapping the execution engine to a Go-based Codex-native implementation.

- ship a compiled Go CLI named `moai-codex`
- keep the MoAI subcommand vocabulary and `.moai/` project structure
- generate compatibility artifacts and project docs locally
- support optional `codex exec` handoff with `--execute`
- stay dependency-light by using the Go standard library only

## Current Scope

Implemented now:

- `moai-codex init`
- `moai-codex update`
- `moai-codex status`
- `moai-codex doctor`
- `moai-codex project`
- `moai-codex plan`
- `moai-codex run`
- `moai-codex sync`
- `moai-codex review`
- `moai-codex coverage`
- `moai-codex clean`
- `moai-codex fix`
- `moai-codex loop`
- `moai-codex codemaps`
- `moai-codex cc`
- `moai-codex cg`
- `moai-codex glm`
- `moai-codex worktree list|status|new|remove`

Behavior model:

- `init` and `update` manage a Codex-compatible scaffold and `.moai/manifest.json`
- workflow commands create MoAI-shaped artifacts under `.moai/project/` and `.moai/specs/`
- `--execute` additionally runs `codex exec` with a prepared compatibility prompt
- runtime mode commands update `.moai` state so status output reflects the active mode

## Requirements

- Go 1.26.1+
- Git
- Codex CLI available on `PATH` for `--execute`

## Build

```bash
make build
./bin/moai-codex version
```

Without `make`:

```bash
/home/pi/.local/go/bin/go build -o bin/moai-codex ./cmd/moai-codex
```

## Usage

Run directly:

```bash
./bin/moai-codex init .
./bin/moai-codex status --json
./bin/moai-codex plan "Add auth compatibility"
./bin/moai-codex run SPEC-ADD-AUTH-COMPATIBILITY-001
./bin/moai-codex review
./bin/moai-codex codemaps
```

To hand a workflow to Codex non-interactively:

```bash
./bin/moai-codex plan "Add auth compatibility" --execute
./bin/moai-codex run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## Generated Files

`moai-codex init` creates:

- `AGENTS.md`
- `.moai/config/sections/project.yaml`
- `.moai/config/sections/workflow.yaml`
- `.moai/config/sections/quality.yaml`
- `.moai/config/sections/llm.yaml`
- `.moai/config/sections/compatibility.yaml`
- `.moai/project/product.md`
- `.moai/project/structure.md`
- `.moai/project/tech.md`
- `.moai/project/codemaps/overview.md`
- `.moai/state/runtime.json`
- `.moai/manifest.json`

## Testing

```bash
make test
```

The Go test suite covers scaffold creation, update safety, status/doctor behavior, SPEC generation, run artifacts, project/codemap refresh, and runtime mode switching.

## Notes

- This is a compatibility-first reimplementation, not a byte-for-byte port of Claude-specific MoAI internals.
- Home-directory `~/.moai` state is intentionally not treated as a project root.
- Runtime log/state noise under `.omx/` is ignored from git by default.
