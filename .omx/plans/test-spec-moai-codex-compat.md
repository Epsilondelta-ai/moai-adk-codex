# Test Spec: MoAI Codex Compatibility Layer

## Verification Targets

### CLI bootstrap

- `node bin/moai.js init .` succeeds in an empty temp project.
- Re-running `update` preserves managed files and reports applied changes.

### Health checks

- `node bin/moai.js status --json` returns machine-readable scaffold status.
- `node bin/moai.js doctor --json` returns checks for runtime, git, and `.moai` presence.

### Command routing

- All declared commands parse correctly.
- Compatibility commands create or update expected artifacts under `.moai/project/`.

### Regression safety

- User-modified managed files are detected through the manifest and not silently replaced.

### Documentation

- README examples correspond to implemented commands and expected output shape.
