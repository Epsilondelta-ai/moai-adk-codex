# moai-adk-codex

MoAI-ADK compatibility layer for Codex and OMX.

This project preserves the familiar MoAI command surface while swapping the execution engine to Codex-native behavior. The current implementation is intentionally pragmatic:

- keep `moai` command names and `.moai/` project structure
- generate compatibility artifacts and project docs locally
- support optional `codex exec` handoff with `--execute`
- avoid new runtime dependencies by using Node.js built-ins only

## Current Scope

Implemented now:

- `moai init`
- `moai update`
- `moai status`
- `moai doctor`
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
- `moai cc`
- `moai cg`
- `moai glm`
- `moai worktree list|status|new|remove`

Behavior model:

- `init` and `update` manage a Codex-compatible scaffold and `.moai/manifest.json`
- workflow commands create MoAI-shaped artifacts under `.moai/project/` and `.moai/specs/`
- `--execute` additionally runs `codex exec` with a prepared compatibility prompt
- runtime mode commands update `.moai` state so status output reflects the active mode

## Requirements

- Node.js 22+
- Git
- Codex CLI available on `PATH` for `--execute`

## Usage

Run directly:

```bash
node bin/moai.js init .
node bin/moai.js status --json
node bin/moai.js plan "Add auth compatibility"
node bin/moai.js run SPEC-ADD-AUTH-COMPATIBILITY-001
node bin/moai.js review
node bin/moai.js codemaps
```

Or install the local binary into your shell:

```bash
npm link
moai init .
moai plan "Add auth compatibility"
```

To hand a workflow to Codex non-interactively:

```bash
moai plan "Add auth compatibility" --execute
moai run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## Generated Files

`moai init` creates:

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
npm test
```

The test suite uses the built-in Node test runner and covers scaffold creation, update safety, status/doctor behavior, SPEC generation, run artifacts, and runtime mode switching.

## Notes

- This is a compatibility-first reimplementation, not a byte-for-byte port of Claude-specific MoAI internals.
- Home-directory `~/.moai` state is intentionally not treated as a project root.
- Runtime log/state noise under `.omx/` is ignored from git by default.
