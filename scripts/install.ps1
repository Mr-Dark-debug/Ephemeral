# Ephemeral Installer for Windows (PowerShell)

$ErrorActionPreference = "Stop"

$REPO_OWNER = "Mr-Dark-debug"
$REPO_NAME = "Ephemeral"
$BINARY_NAME = "ephemeral.exe"

Write-Host "Starting Ephemeral Installation..." -ForegroundColor Blue

# Detect Arch
$ARCH = "amd64" # Default
if ($env:PROCESSOR_ARCHITECTURE -eq "ARM64") { $ARCH = "arm64" }

Write-Host "Detected: windows/$ARCH" -ForegroundColor Green

# Fetch latest release URL
$API_URL = "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest"
$release = Invoke-RestMethod -Uri $API_URL
$asset = $release.assets | Where-Object { $_.name -like "*windows_$ARCH*" } | Select-Object -First 1

if ($null -eq $asset) {
    Write-Error "Could not find a release for Windows/$ARCH"
    exit 1
}

$DOWNLOAD_URL = $asset.browser_download_url

# Download and Install
$INSTALL_DIR = Join-Path $env:USERPROFILE ".local\bin"
if (!(Test-Path $INSTALL_DIR)) { New-Item -ItemType Directory -Path $INSTALL_DIR }

$TEMP_FILE = Join-Path $env:TEMP "ephemeral.zip"
Write-Host "Downloading from $DOWNLOAD_URL..." -ForegroundColor Blue
Invoke-WebRequest -Uri $DOWNLOAD_URL -OutFile $TEMP_FILE

Write-Host "Installing to $INSTALL_DIR..." -ForegroundColor Green
Expand-Archive -Path $TEMP_FILE -DestinationPath $INSTALL_DIR -Force

# Cleanup zip, keeping only the binary
Remove-Item $TEMP_FILE

# PATH check
$gobin = $INSTALL_DIR
if ($env:PATH -notlike "*$gobin*") {
    Write-Host "Warning: $gobin is not in your PATH." -ForegroundColor Red
    Write-Host "Run this to add it permanently:"
    Write-Host "[Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';$gobin', 'User')"
}

Write-Host "Ephemeral installed successfully!" -ForegroundColor Green
Write-Host "Try running: ephemeral --nick Alice" -ForegroundColor Blue
