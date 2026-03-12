#!/usr/bin/env bash
set -euo pipefail

REPO="Epsilondelta-ai/moai-adk-codex"
BIN_NAME="coai"

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_info() { printf "%b[INFO]%b %s\n" "$BLUE" "$NC" "$1"; }
print_success() { printf "%b[SUCCESS]%b %s\n" "$GREEN" "$NC" "$1"; }
print_warning() { printf "%b[WARNING]%b %s\n" "$YELLOW" "$NC" "$1"; }
print_error() { printf "%b[ERROR]%b %s\n" "$RED" "$NC" "$1" >&2; }

VERSION=""
INSTALL_DIR=""
SOURCE_DIR=""
FORCE_SOURCE=0
TARGET_DIR=""
TMP_DIR=""

cleanup() {
  if [ -n "${TMP_DIR:-}" ] && [ -d "${TMP_DIR:-}" ]; then
    rm -rf "$TMP_DIR"
  fi
}
trap cleanup EXIT

usage() {
  cat <<'EOF'
Usage: install.sh [OPTIONS]

Options:
  --version VERSION      Install a specific release version
  --install-dir DIR      Install to a custom directory
  --source               Skip release download and build from source
  --source-dir DIR       Build from an existing local checkout
  -h, --help             Show help

Examples:
  ./install.sh
  ./install.sh --install-dir "$HOME/.local/bin"
  ./install.sh --source
  ./install.sh --source-dir "$PWD"
EOF
}

parse_args() {
  while [ $# -gt 0 ]; do
    case "$1" in
      --version)
        VERSION="${2:-}"
        shift 2
        ;;
      --install-dir)
        INSTALL_DIR="${2:-}"
        shift 2
        ;;
      --source)
        FORCE_SOURCE=1
        shift
        ;;
      --source-dir)
        SOURCE_DIR="${2:-}"
        FORCE_SOURCE=1
        shift 2
        ;;
      -h|--help)
        usage
        exit 0
        ;;
      *)
        print_error "Unknown option: $1"
        usage
        exit 1
        ;;
    esac
  done
}

detect_platform() {
  local os arch
  os="$(uname -s | tr '[:upper:]' '[:lower:]')"
  arch="$(uname -m)"

  case "$os" in
    linux) OS="linux" ;;
    darwin) OS="darwin" ;;
    *)
      print_error "Unsupported OS: $os"
      exit 1
      ;;
  esac

  case "$arch" in
    i386|i486|i586|i686) ARCH="386" ;;
    x86_64|amd64) ARCH="amd64" ;;
    armv6l) ARCH="armv6" ;;
    armv7l) ARCH="armv7" ;;
    arm64|aarch64) ARCH="arm64" ;;
    ppc64le) ARCH="ppc64le" ;;
    riscv64) ARCH="riscv64" ;;
    s390x) ARCH="s390x" ;;
    *)
      print_error "Unsupported architecture: $arch"
      exit 1
      ;;
  esac

  PLATFORM="${OS}_${ARCH}"
  print_success "Detected platform: $PLATFORM"
}

resolve_install_dir() {
  if [ -n "$INSTALL_DIR" ]; then
    TARGET_DIR="$INSTALL_DIR"
  elif command -v go >/dev/null 2>&1; then
    local gobin gopath
    gobin="$(go env GOBIN 2>/dev/null || true)"
    gopath="$(go env GOPATH 2>/dev/null || true)"
    if [ -n "$gobin" ]; then
      TARGET_DIR="$gobin"
    elif [ -n "$gopath" ]; then
      TARGET_DIR="$gopath/bin"
    else
      TARGET_DIR="$HOME/.local/bin"
    fi
  else
    TARGET_DIR="$HOME/.local/bin"
  fi

  mkdir -p "$TARGET_DIR"
}

get_latest_version() {
  local api_url="https://api.github.com/repos/${REPO}/releases/latest"
  local response=""
  if command -v curl >/dev/null 2>&1; then
    response="$(curl -fsSL "$api_url" 2>/dev/null || true)"
  elif command -v wget >/dev/null 2>&1; then
    response="$(wget -qO- "$api_url" 2>/dev/null || true)"
  fi

  if [ -z "$response" ]; then
    return 1
  fi

  VERSION="$(printf "%s" "$response" | grep -o '"tag_name"[[:space:]]*:[[:space:]]*"[^"]*"' | head -n1 | sed -E 's/.*"([^"]+)"/\1/' | sed 's/^v//')"
  [ -n "$VERSION" ]
}

download_file() {
  local url="$1"
  local dest="$2"
  if command -v curl >/dev/null 2>&1; then
    curl -fsSL "$url" -o "$dest"
  else
    wget -q "$url" -O "$dest"
  fi
}

install_binary() {
  local binary_path="$1"
  local target_path="$TARGET_DIR/$BIN_NAME"
  cp "$binary_path" "$target_path"
  chmod +x "$target_path"
  print_success "Installed to: $target_path"
}

try_release_install() {
  [ -n "$VERSION" ] || get_latest_version || return 1

  TMP_DIR="$(mktemp -d)"
  local archive_name="${BIN_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
  local archive_path="$TMP_DIR/$archive_name"
  local download_url="https://github.com/${REPO}/releases/download/v${VERSION}/${archive_name}"

  print_info "Trying GitHub release install: $download_url"
  if ! download_file "$download_url" "$archive_path"; then
    print_warning "Release download unavailable for version ${VERSION}"
    return 1
  fi

  tar -xzf "$archive_path" -C "$TMP_DIR"
  if [ ! -f "$TMP_DIR/$BIN_NAME" ]; then
    print_warning "Release archive did not contain ${BIN_NAME}"
    return 1
  fi

  install_binary "$TMP_DIR/$BIN_NAME"
  return 0
}

build_from_source() {
  TMP_DIR="$(mktemp -d)"
  local source_root="$SOURCE_DIR"

  if [ -z "$source_root" ]; then
    if ! command -v git >/dev/null 2>&1; then
      print_error "git is required for source fallback"
      exit 1
    fi
    source_root="$TMP_DIR/src"
    print_info "Cloning source from GitHub"
    git clone --depth 1 "https://github.com/${REPO}.git" "$source_root" >/dev/null 2>&1
  fi

  local go_bin=""
  if command -v go >/dev/null 2>&1; then
    go_bin="go"
  elif [ -x "$HOME/.local/go/bin/go" ]; then
    go_bin="$HOME/.local/go/bin/go"
  else
    print_error "Go toolchain not found for source fallback"
    exit 1
  fi

  print_info "Building ${BIN_NAME} from source"
  (
    cd "$source_root"
    "$go_bin" build -o "$TMP_DIR/$BIN_NAME" ./cmd/coai
  )
  install_binary "$TMP_DIR/$BIN_NAME"
}

verify_installation() {
  local target_path="$TARGET_DIR/$BIN_NAME"
  if [ -x "$target_path" ]; then
    print_success "Installed successfully"
    "$target_path" version || true
  fi

  case ":$PATH:" in
    *":$TARGET_DIR:"*) ;;
    *)
      print_warning "Add ${TARGET_DIR} to PATH to use '${BIN_NAME}' globally"
      ;;
  esac
}

main() {
  parse_args "$@"
  detect_platform
  resolve_install_dir

  if [ "$FORCE_SOURCE" -eq 1 ]; then
    build_from_source
    verify_installation
    return
  fi

  if try_release_install; then
    verify_installation
    return
  fi

  print_warning "Falling back to source build"
  build_from_source
  verify_installation
}

main "$@"
