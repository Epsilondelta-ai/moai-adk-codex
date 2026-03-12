# coai

coai は、元の [MoAI-ADK](https://github.com/modu-ai/moai-adk) に着想を得た、Codex 向けの再実装 / 互換レイヤーです。

[English](./README.md) · [한국어](./README.ko.md) · [日本語](./README.ja.md) · [中文](./README.zh.md)

## 概要

- `coai` という名前の Go CLI
- MoAI 風のサブコマンドを維持した Codex 向けワークフロー面
- `.coai/` を正規の状態ディレクトリとして使用し、従来の `.moai/` は読み取り互換を維持
- 標準ライブラリのみで実装し、必要に応じて `codex exec` に委譲可能

## 現在のコマンド

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

## リリース対応ターゲット

release asset は次のターゲット向けに生成されます。

- Linux: `386`, `amd64`, `armv6`, `armv7`, `arm64`, `ppc64le`, `riscv64`, `s390x`
- macOS: `amd64`, `arm64`
- Windows: `386`, `amd64`, `arm64`

## インストール

元の `moai-adk` に近い形のインストールスクリプトを提供しています。

### macOS / Linux / WSL

```bash
curl -fsSL https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.sh | bash
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.ps1 | iex
```

### Windows Batch

`install.bat` をダウンロードして実行します。

メモ:

- まず GitHub Release のバイナリを試します。
- 一致する release asset があれば、利用者は Go をインストールする必要はありません。
- 一致する release asset がなければ、ソースビルドへ自動的にフォールバックし、その場合は Go が必要です。
- ローカル確認には `install.sh --source-dir "$PWD"` または `install.ps1 -SourceDir $PWD` を使えます。

release asset は、このリポジトリに `v0.3.0` のようなタグを push すると公開されます。

## ビルド

```bash
make build
./bin/coai version
```

直接ビルド:

```bash
/home/pi/.local/go/bin/go build -o bin/coai ./cmd/coai
```

## 使い方

```bash
./bin/coai init .
./bin/coai status --json
./bin/coai plan "認証機能を追加"
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001
./bin/coai review
./bin/coai codemaps
```

Codex に直接渡す場合:

```bash
./bin/coai plan "認証機能を追加" --execute
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## 生成ファイル

`coai init` は次を生成します。

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

## テスト

```bash
make test
```

## 補足

- 元のプロジェクトは [MoAI-ADK](https://github.com/modu-ai/moai-adk) です。
- coai は `.coai/` を正規の状態ディレクトリとして使います。
- 既存の `.moai/` ディレクトリは互換性のため読み取りのみサポートします。
