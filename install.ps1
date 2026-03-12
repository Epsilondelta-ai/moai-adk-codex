$ErrorActionPreference = "Stop"

$Repo = "Epsilondelta-ai/moai-adk-codex"
$BinName = "coai"
$Version = ""
$InstallDir = ""
$SourceDir = ""
$UseSource = $false

function Write-Info($Message) { Write-Host "[INFO] $Message" -ForegroundColor Cyan }
function Write-Success($Message) { Write-Host "[SUCCESS] $Message" -ForegroundColor Green }
function Write-WarningMsg($Message) { Write-Host "[WARNING] $Message" -ForegroundColor Yellow }
function Write-ErrorMsg($Message) { Write-Host "[ERROR] $Message" -ForegroundColor Red }

function Show-Usage {
  @"
Usage: install.ps1 [-Version <version>] [-InstallDir <dir>] [-Source] [-SourceDir <dir>]

Examples:
  powershell -ExecutionPolicy Bypass -File .\install.ps1
  powershell -ExecutionPolicy Bypass -File .\install.ps1 -Source
  powershell -ExecutionPolicy Bypass -File .\install.ps1 -SourceDir $PWD
"@
}

function Resolve-InstallDir {
  if ($InstallDir) { return $InstallDir }

  $go = Get-Command go -ErrorAction SilentlyContinue
  if ($go) {
    $gobin = & go env GOBIN 2>$null
    if ($gobin) { return $gobin.Trim() }
    $gopath = & go env GOPATH 2>$null
    if ($gopath) { return (Join-Path $gopath.Trim() "bin") }
  }

  return (Join-Path $HOME ".local\bin")
}

function Get-Platform {
  $arch = if ($env:ARCHITEW6432) { $env:ARCHITEW6432 } else { $env:PROCESSOR_ARCHITECTURE }
  switch ($arch.ToUpper()) {
    "AMD64" { return "windows_amd64" }
    "X86"   { return "windows_386" }
    "ARM64" { return "windows_arm64" }
    default { throw "Unsupported architecture: $arch" }
  }
}

function Get-LatestVersion {
  $response = Invoke-RestMethod -Uri "https://api.github.com/repos/$Repo/releases/latest" -Method Get -ErrorAction SilentlyContinue
  if (-not $response) { return $null }
  return $response.tag_name.TrimStart("v")
}

function Try-ReleaseInstall($TargetDir, $Platform) {
  if (-not $Version) {
    $script:Version = Get-LatestVersion
  }
  if (-not $Version) { return $false }

  $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("coai-install-" + [guid]::NewGuid().ToString("N"))
  New-Item -ItemType Directory -Path $tmp | Out-Null
  try {
    $archiveName = "$BinName" + "_" + "$Version" + "_" + "$Platform.zip"
    $downloadUrl = "https://github.com/$Repo/releases/download/v$Version/$archiveName"
    $archivePath = Join-Path $tmp $archiveName
    Write-Info "Trying GitHub release install: $downloadUrl"
    Invoke-WebRequest -Uri $downloadUrl -OutFile $archivePath -ErrorAction Stop | Out-Null
    Expand-Archive -Path $archivePath -DestinationPath $tmp -Force

    $binaryPath = Join-Path $tmp "$BinName.exe"
    if (-not (Test-Path $binaryPath)) {
      Write-WarningMsg "Release archive did not contain $BinName.exe"
      return $false
    }

    New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
    Copy-Item $binaryPath (Join-Path $TargetDir "$BinName.exe") -Force
    Write-Success "Installed to $(Join-Path $TargetDir "$BinName.exe")"
    return $true
  } catch {
    Write-WarningMsg "Release download unavailable: $_"
    return $false
  } finally {
    Remove-Item $tmp -Recurse -Force -ErrorAction SilentlyContinue
  }
}

function Build-FromSource($TargetDir) {
  $goBin = $null
  $goCmd = Get-Command go -ErrorAction SilentlyContinue
  if ($goCmd) {
    $goBin = "go"
  } elseif (Test-Path "$HOME\.local\go\bin\go.exe") {
    $goBin = "$HOME\.local\go\bin\go.exe"
  } elseif (Test-Path "$HOME\.local\go\bin\go") {
    $goBin = "$HOME\.local\go\bin\go"
  }
  if (-not $goBin) { throw "Go toolchain not found for source fallback" }

  $tmp = Join-Path ([System.IO.Path]::GetTempPath()) ("coai-src-" + [guid]::NewGuid().ToString("N"))
  $src = $SourceDir
  if (-not $src) {
    $gitCmd = Get-Command git -ErrorAction SilentlyContinue
    if (-not $gitCmd) { throw "git is required for source fallback" }
    Write-Info "Cloning source from GitHub"
    & git clone --depth 1 "https://github.com/$Repo.git" $tmp | Out-Null
    $src = $tmp
  }

  New-Item -ItemType Directory -Path $TargetDir -Force | Out-Null
  $target = Join-Path $TargetDir "$BinName.exe"
  Write-Info "Building $BinName from source"
  Push-Location $src
  try {
    & $goBin build -o $target .\cmd\coai
  } finally {
    Pop-Location
  }
  Write-Success "Installed to $target"
}

param(
  [string]$Version,
  [string]$InstallDir,
  [switch]$Source,
  [string]$SourceDir,
  [switch]$Help
)

if ($Help) {
  Show-Usage
  exit 0
}

$UseSource = $Source.IsPresent -or [bool]$SourceDir
$targetDir = Resolve-InstallDir
$platform = Get-Platform

if (-not $UseSource) {
  if (Try-ReleaseInstall $targetDir $platform) {
    exit 0
  }
}

Write-WarningMsg "Falling back to source build"
Build-FromSource $targetDir
