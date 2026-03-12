const path = require("node:path");

const { ensureProjectScaffold, readProjectStatus, runDoctor, setRuntimeMode } = require("./core/project");
const { renderJson, renderText, slugify, timestampUtc } = require("./core/utils");
const { createSpec, createWorkflowArtifact, maybeExecuteWithCodex, runWorktreeCommand } = require("./workflows/commands");

const HELP_TEXT = `MoAI Codex Compatibility CLI

Usage:
  moai <command> [args]

Core commands:
  init [path]        Initialize a Codex-compatible MoAI scaffold
  update [path]      Reapply managed scaffold files
  status             Show scaffold and mode status
  doctor             Check runtime and scaffold health

Workflow commands:
  project [summary]
  plan <description>
  run <spec-id|description>
  sync [spec-id]
  review
  coverage
  clean
  fix
  loop
  codemaps

Mode commands:
  cc | cg | glm

Git worktree:
  worktree list
  worktree status
  worktree new <name>
  worktree remove <path>

Flags:
  --json       Render machine-readable output where supported
  --execute    For workflow commands, also invoke codex exec
`;

async function main(argv, io = { stdout: process.stdout, stderr: process.stderr, cwd: process.cwd() }) {
  const { command, args, flags } = parseArgs(argv);

  switch (command) {
    case "help":
    case "--help":
    case "-h":
      io.stdout.write(`${HELP_TEXT}\n`);
      return;
    case "version":
    case "--version":
    case "-V":
      io.stdout.write("moai-adk-codex 0.1.0\n");
      return;
    case "init": {
      const targetDir = path.resolve(io.cwd, args[0] || ".");
      const result = ensureProjectScaffold(targetDir, { forceUpdate: false });
      writeOutput(io, flags.json, result);
      return;
    }
    case "update": {
      const targetDir = path.resolve(io.cwd, args[0] || ".");
      const result = ensureProjectScaffold(targetDir, { forceUpdate: true });
      writeOutput(io, flags.json, result);
      return;
    }
    case "status": {
      const status = readProjectStatus(io.cwd);
      writeOutput(io, flags.json, status);
      return;
    }
    case "doctor": {
      const report = runDoctor(io.cwd);
      writeOutput(io, flags.json, report);
      return;
    }
    case "cc":
    case "cg":
    case "glm": {
      const result = setRuntimeMode(io.cwd, command);
      writeOutput(io, flags.json, result);
      return;
    }
    case "project":
    case "sync":
    case "review":
    case "coverage":
    case "clean":
    case "fix":
    case "loop":
    case "codemaps": {
      ensureProjectScaffold(io.cwd, { forceUpdate: false });
      const artifact = createWorkflowArtifact(io.cwd, command, args.join(" ").trim());
      if (flags.execute) {
        artifact.execution = maybeExecuteWithCodex(io.cwd, command, artifact.prompt);
      }
      writeOutput(io, flags.json, artifact);
      return;
    }
    case "plan": {
      ensureProjectScaffold(io.cwd, { forceUpdate: false });
      if (args.length === 0) {
        throw new Error("plan requires a description");
      }
      const spec = createSpec(io.cwd, args.join(" "));
      if (flags.execute) {
        spec.execution = maybeExecuteWithCodex(io.cwd, "plan", spec.prompt);
      }
      writeOutput(io, flags.json, spec);
      return;
    }
    case "run": {
      ensureProjectScaffold(io.cwd, { forceUpdate: false });
      if (args.length === 0) {
        throw new Error("run requires a SPEC id or description");
      }
      const raw = args.join(" ").trim();
      const payload = createWorkflowArtifact(io.cwd, "run", raw, {
        specId: raw.startsWith("SPEC-") ? raw : `SPEC-${slugify(raw).toUpperCase()}-001`
      });
      if (flags.execute) {
        payload.execution = maybeExecuteWithCodex(io.cwd, "run", payload.prompt);
      }
      writeOutput(io, flags.json, payload);
      return;
    }
    case "worktree": {
      const result = runWorktreeCommand(io.cwd, args);
      writeOutput(io, flags.json, result);
      return;
    }
    case "":
      io.stdout.write(`${HELP_TEXT}\n`);
      return;
    default:
      throw new Error(`unknown command: ${command}`);
  }
}

function writeOutput(io, asJson, payload) {
  const output = asJson ? renderJson(payload) : renderText(payload);
  io.stdout.write(`${output}\n`);
}

function parseArgs(argv) {
  const flags = { json: false, execute: false };
  const remaining = [];

  for (const value of argv) {
    if (value === "--json") {
      flags.json = true;
      continue;
    }
    if (value === "--execute") {
      flags.execute = true;
      continue;
    }
    remaining.push(value);
  }

  return {
    command: remaining[0] || "",
    args: remaining.slice(1),
    flags
  };
}

module.exports = {
  HELP_TEXT,
  main,
  parseArgs,
  timestampUtc
};
