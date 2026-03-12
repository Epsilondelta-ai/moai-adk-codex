# coai

coai 是一个面向 Codex 的重实现 / 兼容层，其灵感来自原始项目 [MoAI-ADK](https://github.com/modu-ai/moai-adk)。

[English](./README.md) · [한국어](./README.ko.md) · [日本語](./README.ja.md) · [中文](./README.zh.md)

## 简介

- 一个名为 `coai` 的 Go CLI
- 保留 MoAI 风格子命令的 Codex 工作流界面
- 使用 `.coai/` 作为规范状态目录，同时保留对旧 `.moai/` 的读取兼容
- 仅使用标准库实现，并可在需要时委托给 `codex exec`

## 当前命令

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

## 安装

仓库提供了与原始 `moai-adk` 类似的安装入口。

### macOS / Linux / WSL

```bash
curl -fsSL https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.sh | bash
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.ps1 | iex
```

### Windows Batch

下载并运行 `install.bat`。

说明：

- 安装器会优先尝试 GitHub Release 二进制包。
- 如果存在匹配的发布资产，用户不需要安装 Go。
- 如果没有匹配的发布资产，会自动回退到源码构建，此时需要 Go。
- 本地测试可使用 `install.sh --source-dir "$PWD"` 或 `install.ps1 -SourceDir $PWD`。

发布资产会在向此仓库推送 `v0.3.0` 这类标签时自动生成。

## 构建

```bash
make build
./bin/coai version
```

直接构建：

```bash
/home/pi/.local/go/bin/go build -o bin/coai ./cmd/coai
```

## 用法

```bash
./bin/coai init .
./bin/coai status --json
./bin/coai plan "Add auth compatibility"
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001
./bin/coai review
./bin/coai codemaps
```

若要直接交给 Codex：

```bash
./bin/coai plan "Add auth compatibility" --execute
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## 生成文件

`coai init` 会生成：

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

## 测试

```bash
make test
```

## 备注

- 原始项目是 [MoAI-ADK](https://github.com/modu-ai/moai-adk)。
- coai 使用 `.coai/` 作为规范状态目录。
- 旧的 `.moai/` 目录仍保留只读兼容能力。
