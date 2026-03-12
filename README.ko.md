# coai

coai는 원본 [MoAI-ADK](https://github.com/modu-ai/moai-adk)에서 영감을 받은 Codex 중심 재구현/호환 레이어입니다.

[English](./README.md) · [한국어](./README.ko.md) · [日本語](./README.ja.md) · [中文](./README.zh.md)

## 개요

- `coai`라는 이름의 Go CLI
- MoAI 스타일 서브커맨드를 유지한 Codex 중심 워크플로우 표면
- `.coai/`를 기본 상태 디렉터리로 사용하고, 기존 `.moai/`는 읽기 호환 지원
- 표준 라이브러리만 사용하며 필요 시 `codex exec`로 위임 가능

## 현재 지원 명령

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

## 릴리스 지원 타깃

릴리스 자산은 다음 타깃을 대상으로 준비됩니다.

- Linux: `386`, `amd64`, `armv6`, `armv7`, `arm64`, `ppc64le`, `riscv64`, `s390x`
- macOS: `amd64`, `arm64`
- Windows: `386`, `amd64`, `arm64`

## 설치

원본 `moai-adk`와 비슷한 방식의 설치 스크립트를 제공합니다.

### macOS / Linux / WSL

```bash
curl -fsSL https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.sh | bash
```

### Windows PowerShell

```powershell
irm https://raw.githubusercontent.com/Epsilondelta-ai/moai-adk-codex/main/install.ps1 | iex
```

### Windows Batch

`install.bat`를 내려받아 실행하면 됩니다.

참고:

- 먼저 GitHub Release 바이너리를 시도합니다.
- 일치하는 릴리스 자산이 있으면 사용자는 Go를 설치하지 않아도 됩니다.
- 일치하는 릴리스 자산이 없으면 소스 빌드로 자동 fallback 하며, 이 경우 Go가 필요합니다.
- 로컬 테스트는 `install.sh --source-dir "$PWD"` 또는 `install.ps1 -SourceDir $PWD`를 사용하면 됩니다.

릴리스 자산은 이 저장소에 `v0.4.0` 같은 태그를 푸시하면 발행됩니다.

## 빌드

```bash
make build
./bin/coai version
```

직접 빌드:

```bash
/home/pi/.local/go/bin/go build -o bin/coai ./cmd/coai
```

## 사용 예시

```bash
./bin/coai init .
./bin/coai status --json
./bin/coai plan "인증 기능 추가"
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001
./bin/coai review
./bin/coai codemaps
```

Codex에 바로 넘기려면:

```bash
./bin/coai plan "인증 기능 추가" --execute
./bin/coai run SPEC-ADD-AUTH-COMPATIBILITY-001 --execute
```

## 생성 파일

`coai init`는 다음을 생성합니다.

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

## 테스트

```bash
make test
```

## 비고

- 원본 프로젝트는 [MoAI-ADK](https://github.com/modu-ai/moai-adk)입니다.
- coai는 `.coai/`를 기본 상태 디렉터리로 사용합니다.
- 기존 `.moai/` 디렉터리는 호환을 위해 읽기만 지원합니다.
