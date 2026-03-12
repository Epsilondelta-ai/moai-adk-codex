const crypto = require("node:crypto");
const fs = require("node:fs");
const path = require("node:path");

function ensureDir(dirPath) {
  fs.mkdirSync(dirPath, { recursive: true });
}

function readJson(filePath, fallback) {
  try {
    return JSON.parse(fs.readFileSync(filePath, "utf8"));
  } catch (error) {
    return fallback;
  }
}

function writeJson(filePath, data) {
  ensureDir(path.dirname(filePath));
  fs.writeFileSync(filePath, `${JSON.stringify(data, null, 2)}\n`, "utf8");
}

function sha256(content) {
  return crypto.createHash("sha256").update(content).digest("hex");
}

function timestampUtc(date = new Date()) {
  return date.toISOString().replaceAll(":", "").replace(/\.\d{3}Z$/, "Z");
}

function slugify(value) {
  return value
    .toLowerCase()
    .replace(/[^a-z0-9]+/g, "-")
    .replace(/^-+|-+$/g, "")
    .replace(/-{2,}/g, "-") || "task";
}

function readSimpleYaml(filePath) {
  const result = {};
  if (!fs.existsSync(filePath)) {
    return result;
  }

  const lines = fs.readFileSync(filePath, "utf8").split(/\r?\n/);
  const stack = [{ indent: -1, target: result }];

  for (const line of lines) {
    if (!line.trim() || line.trimStart().startsWith("#")) {
      continue;
    }

    const indent = line.match(/^ */)[0].length;
    const trimmed = line.trim();
    const [rawKey, ...rest] = trimmed.split(":");
    const key = rawKey.trim();
    const rawValue = rest.join(":").trim();

    while (stack.length > 1 && indent <= stack[stack.length - 1].indent) {
      stack.pop();
    }

    const current = stack[stack.length - 1].target;

    if (rawValue === "") {
      current[key] = {};
      stack.push({ indent, target: current[key] });
      continue;
    }

    current[key] = normalizeYamlScalar(rawValue);
  }

  return result;
}

function normalizeYamlScalar(value) {
  if (value === "true") {
    return true;
  }
  if (value === "false") {
    return false;
  }
  if (/^-?\d+$/.test(value)) {
    return Number(value);
  }
  return value.replace(/^["']|["']$/g, "");
}

function renderJson(payload) {
  return JSON.stringify(payload, null, 2);
}

function renderText(payload) {
  if (typeof payload === "string") {
    return payload;
  }
  if (Array.isArray(payload)) {
    return payload.map(renderText).join("\n");
  }
  return Object.entries(payload)
    .map(([key, value]) => `${key}: ${Array.isArray(value) || typeof value === "object" ? JSON.stringify(value) : value}`)
    .join("\n");
}

module.exports = {
  ensureDir,
  normalizeYamlScalar,
  readJson,
  readSimpleYaml,
  renderJson,
  renderText,
  sha256,
  slugify,
  timestampUtc,
  writeJson
};
