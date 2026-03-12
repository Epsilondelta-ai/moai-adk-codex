const fs = require("node:fs");
const os = require("node:os");
const path = require("node:path");
const childProcess = require("node:child_process");

const { buildTemplates } = require("./templates");
const { ensureDir, readJson, readSimpleYaml, sha256, timestampUtc, writeJson } = require("./utils");

function ensureProjectScaffold(projectRoot, options = {}) {
  const root = path.resolve(projectRoot);
  ensureDir(root);

  const context = {
    projectName: path.basename(root),
    initializedAt: new Date().toISOString()
  };
  const manifestPath = path.join(root, ".moai", "manifest.json");
  const manifest = readJson(manifestPath, { files: {} });
  const created = [];
  const updated = [];
  const skipped = [];

  for (const template of buildTemplates(context)) {
    const filePath = path.join(root, template.path);
    const contentHash = sha256(template.content);
    const existing = manifest.files[template.path];

    ensureDir(path.dirname(filePath));

    if (!fs.existsSync(filePath)) {
      fs.writeFileSync(filePath, template.content, "utf8");
      manifest.files[template.path] = { managed: true, hash: contentHash };
      created.push(template.path);
      continue;
    }

    const currentContent = fs.readFileSync(filePath, "utf8");
    const currentHash = sha256(currentContent);
    const isUnmodifiedManaged = existing && existing.hash === currentHash;

    if (options.forceUpdate && isUnmodifiedManaged) {
      fs.writeFileSync(filePath, template.content, "utf8");
      manifest.files[template.path] = { managed: true, hash: contentHash };
      updated.push(template.path);
      continue;
    }

    if (!existing) {
      manifest.files[template.path] = { managed: false, hash: currentHash };
      skipped.push(template.path);
      continue;
    }

    if (existing.hash !== currentHash) {
      skipped.push(template.path);
      continue;
    }

    manifest.files[template.path] = { managed: true, hash: contentHash };
  }

  writeJson(manifestPath, manifest);

  return {
    command: options.forceUpdate ? "update" : "init",
    projectRoot: root,
    created,
    updated,
    skipped,
    manifest: path.relative(root, manifestPath)
  };
}

function readProjectStatus(projectRoot) {
  const root = findProjectRoot(projectRoot);
  const runtimePath = path.join(root, ".moai", "state", "runtime.json");
  const runtime = readJson(runtimePath, { currentRuntimeMode: "cg", currentSpec: null, lastCommand: null });
  const projectConfig = readSimpleYaml(path.join(root, ".moai", "config", "sections", "project.yaml"));
  const quality = readSimpleYaml(path.join(root, ".moai", "config", "sections", "quality.yaml"));
  const manifest = readJson(path.join(root, ".moai", "manifest.json"), { files: {} });

  return {
    projectRoot: root,
    projectName: projectConfig.project?.name || path.basename(root),
    runtimeMode: runtime.currentRuntimeMode || "cg",
    currentSpec: runtime.currentSpec,
    lastCommand: runtime.lastCommand,
    developmentMode: quality.constitution?.development_mode || "tdd",
    managedFiles: Object.keys(manifest.files).length,
    initialized: fs.existsSync(path.join(root, ".moai"))
  };
}

function setRuntimeMode(projectRoot, mode) {
  const root = findProjectRoot(projectRoot);
  const runtimePath = path.join(root, ".moai", "state", "runtime.json");
  const runtime = readJson(runtimePath, {});
  runtime.currentRuntimeMode = mode;
  runtime.lastCommand = mode;
  writeJson(runtimePath, runtime);

  const llmPath = path.join(root, ".moai", "config", "sections", "llm.yaml");
  fs.writeFileSync(
    llmPath,
    `llm:\n  current_runtime_mode: ${mode}\n  provider: codex\n`,
    "utf8"
  );

  return {
    command: mode,
    runtimeMode: mode,
    projectRoot: root
  };
}

function runDoctor(projectRoot) {
  const root = path.resolve(projectRoot);
  const checks = [];

  checks.push(checkBinary("node", "Node.js runtime"));
  checks.push(checkBinary("git", "Git"));
  checks.push(checkBinary("codex", "Codex CLI"));

  const moaiDir = path.join(root, ".moai");
  checks.push({
    name: ".moai scaffold",
    ok: fs.existsSync(moaiDir),
    details: fs.existsSync(moaiDir) ? "present" : "missing"
  });

  return {
    projectRoot: root,
    ok: checks.every((entry) => entry.ok),
    checks
  };
}

function checkBinary(binary, name) {
  try {
    childProcess.execFileSync("which", [binary], { stdio: "ignore" });
    return { name, ok: true, details: "available" };
  } catch (error) {
    return { name, ok: false, details: "missing" };
  }
}

function findProjectRoot(startDir) {
  const initial = path.resolve(startDir);
  let current = initial;
  const homeDir = path.resolve(os.homedir());

  while (true) {
    const hasMoai = fs.existsSync(path.join(current, ".moai"));
    if (hasMoai && current !== homeDir) {
      return current;
    }

    const parent = path.dirname(current);
    if (parent === current) {
      return initial;
    }
    current = parent;
  }
}

function updateRuntime(projectRoot, patch) {
  const root = findProjectRoot(projectRoot);
  const runtimePath = path.join(root, ".moai", "state", "runtime.json");
  const runtime = readJson(runtimePath, {});
  writeJson(runtimePath, { ...runtime, ...patch, updatedAt: new Date().toISOString() });
}

module.exports = {
  ensureProjectScaffold,
  findProjectRoot,
  readProjectStatus,
  runDoctor,
  setRuntimeMode,
  updateRuntime
};
