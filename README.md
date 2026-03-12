# coai

MoAI-ADK compatibility layer for Codex and OMX, exposed as the `coai` command.

This project preserves the familiar MoAI command surface while swapping the execution engine to a Go-based Codex-native implementation.

- ship a compiled Go CLI named `coai`
- keep the MoAI subcommand vocabulary while using `.coai/` as the canonical project structure
- generate compatibility artifacts and project docs locally
- support optional `codex exec` handoff with `--execute`
- stay dependency-light by using the Go standard library only

## Current Scope

Implemented now:

- `coai init`
- `coai update`
- `coai status`
- `coai doctor`
- `coai project`
- `coai plan`
- `coai run`
- `coai sync`
- `coai review`
- `coai coverage`
- `coai clean`
- `coai fix`
- `coai loop`
- `coai codemaps`
- `coai cc`
- `coai cg`
- `coai glm`
- `coai worktree list|status|new|remove`

Behavior model:

- `init` and `update` manage a Codex-compatible scaffold and `.coai/manifest.json`
- workflow commands create MoAI-shaped artifacts under `.coai/project/` and `.coai/specs/`
- `--execute` additionally runs `codex exec` with a prepared compatibility prompt
- runtime mode commands update `.coai` state so status output reflects the active mode
- legacy `.moai/` directories are still read for compatibility

## Requirements

- Go 1.26.1+
- Git
- Codex CLI available on `PATH` for `--execute`

## Install

Upstream-style install entrypoints are provided at the repo root:

### macOS / Linux / WSL

```bash
curl -fsSL https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.sh | bash
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.ps1 | iex
```

### Windows Batch

Download and run `install.bat`.

Notes:

- The installer tries GitHub Releases first.
- If a matching release asset is not available, it falls back to building from source.
- For local testing or local builds, use `install.sh --source-dir "$PWD"` or `install.ps1 -SourceDir $PWD`.

## Build

```bash
make build
./bin/coai version
```

Without `make`:

```bash
/home/pi/.local/go/bin/go build -o bin/coai ./cmd/coai
```

## Usage

Run directly:

```bash
./bin/coai init .
./bin/coai status --json
./bin/coai plan "Add auth compatibility"
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001
./bin/coai review
./bin/coai codemaps
```

To hand a workflow to Codex non-interactively:

```bash
./bin/coai plan "Add auth compatibility" --execute
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## Generated Files

`coai init` creates:

- `AGENTS.md`
- `.coai/config/sections/project.yaml`
- `.coai/config/sections/workflow.yaml`
- `.coai/config/sections/quality.yaml`
- `.coai/config/sections/llm.yaml`
- `.coai/config/sections/compatibility.yaml`
- `.coai/project/product.md`
- `.coai/project/structure.md`
- `.coai/project/tech.md`
- `.coai/project/codemaps/overview.md`
- `.coai/state/runtime.json`
- `.coai/manifest.json`

## Testing

```bash
make test
```

The Go test suite covers scaffold creation, update safety, status/doctor behavior, SPEC generation, run artifacts, project/codemap refresh, and runtime mode switching.

## Notes

- This is a compatibility-first reimplementation, not a byte-for-byte port of Claude-specific MoAI internals.
- Home-directory `~/.coai` and legacy `~/.moai` state are intentionally not treated as project roots.
- Runtime log/state noise under `.omx/` is ignored from git by default.
