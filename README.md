# coai

coai is a Codex-oriented reimplementation and compatibility layer inspired by the original [MoAI-ADK](https://github.com/modu-ai/moai-adk).

[English](./README.md) · [한국어](./README.ko.md) · [日本語](./README.ja.md) · [中文](./README.zh.md)

## What It Is

- a compiled Go CLI named `coai`
- a Codex-focused workflow surface that keeps MoAI-style subcommands
- a `.coai/` project state layout with legacy `.moai/` read fallback
- a standard-library-only implementation with optional `codex exec` handoff

## Current Commands

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

## Supported Release Targets

Release assets are prepared for these targets:

- Linux: `386`, `amd64`, `armv6`, `armv7`, `arm64`, `ppc64le`, `riscv64`, `s390x`
- macOS: `amd64`, `arm64`
- Windows: `386`, `amd64`, `arm64`

## Install

Upstream-style installer entrypoints are available in this repository.

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
- If a matching release asset exists, users do not need Go installed.
- If no matching release asset exists, it falls back to building from source, which requires Go.
- For local testing, use `install.sh --source-dir "$PWD"` or `install.ps1 -SourceDir $PWD`.

Release assets are published from this repository by pushing a tag like `v0.3.0`.

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

## Notes

- The original project is [MoAI-ADK](https://github.com/modu-ai/moai-adk).
- coai uses `.coai/` as the canonical state directory.
- Legacy `.moai/` directories are still read for compatibility.
