const test = require("node:test");
const assert = require("node:assert/strict");
const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const childProcess = require("node:child_process");

const repoRoot = path.resolve(__dirname, "..");
const cliPath = path.join(repoRoot, "bin", "moai.js");

function runCli(args, cwd) {
  return childProcess.execFileSync("node", [cliPath, ...args], {
    cwd,
    encoding: "utf8"
  });
}

function makeTempDir() {
  return fs.mkdtempSync(path.join(os.tmpdir(), "moai-codex-"));
}

test("init creates scaffold and manifest", () => {
  const dir = makeTempDir();
  const output = runCli(["init", ".", "--json"], dir);
  const payload = JSON.parse(output);

  assert.equal(payload.command, "init");
  assert.ok(fs.existsSync(path.join(dir, ".moai", "manifest.json")));
  assert.ok(fs.existsSync(path.join(dir, "AGENTS.md")));
});

test("status reflects initialized scaffold", () => {
  const dir = makeTempDir();
  runCli(["init", "."], dir);
  const output = runCli(["status", "--json"], dir);
  const payload = JSON.parse(output);

  assert.equal(payload.initialized, true);
  assert.equal(payload.runtimeMode, "cg");
  assert.equal(payload.developmentMode, "tdd");
});

test("update preserves user-modified files", () => {
  const dir = makeTempDir();
  runCli(["init", ".", "--json"], dir);

  const agentsPath = path.join(dir, "AGENTS.md");
  fs.writeFileSync(agentsPath, "# custom\n", "utf8");

  const output = runCli(["update", ".", "--json"], dir);
  const payload = JSON.parse(output);

  assert.ok(payload.skipped.includes("AGENTS.md"));
  assert.equal(fs.readFileSync(agentsPath, "utf8"), "# custom\n");
});

test("plan creates a spec artifact", () => {
  const dir = makeTempDir();
  runCli(["init", "."], dir);
  const output = runCli(["plan", "Add compatibility routing", "--json"], dir);
  const payload = JSON.parse(output);

  assert.match(payload.specId, /^SPEC-/);
  assert.ok(fs.existsSync(path.join(dir, payload.specPath)));
});

test("run creates a workflow artifact and stores spec", () => {
  const dir = makeTempDir();
  runCli(["init", "."], dir);
  const output = runCli(["run", "SPEC-TEST-001", "--json"], dir);
  const payload = JSON.parse(output);

  assert.equal(payload.specId, "SPEC-TEST-001");
  assert.ok(fs.existsSync(path.join(dir, payload.artifactPath)));
});

test("mode switch updates runtime state", () => {
  const dir = makeTempDir();
  runCli(["init", "."], dir);
  runCli(["glm"], dir);
  const status = JSON.parse(runCli(["status", "--json"], dir));

  assert.equal(status.runtimeMode, "glm");
});

test("doctor reports scaffold availability", () => {
  const dir = makeTempDir();
  runCli(["init", "."], dir);
  const output = runCli(["doctor", "--json"], dir);
  const payload = JSON.parse(output);

  assert.equal(payload.ok, true);
  assert.ok(payload.checks.some((entry) => entry.name === ".moai scaffold" && entry.ok));
});

test("status does not treat home directory .moai as project root", () => {
  const dir = makeTempDir();
  const output = runCli(["status", "--json"], dir);
  const payload = JSON.parse(output);

  assert.equal(payload.projectRoot, dir);
  assert.equal(payload.initialized, false);
});
