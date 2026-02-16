# Ephemeral Installer for Windows (PowerShell)

$ErrorActionPreference = "Stop"

Write-Host "Starting Ephemeral Installation..." -ForegroundColor Blue

# Check for Go
if (!(Get-Command go -ErrorAction SilentlyContinue)) {
    Write-Host "Error: Go is not installed." -ForegroundColor Red
    Write-Host "Please install Go (1.23+) from https://golang.org/dl/"
    exit 1
}

# Install
if (Test-Path "cmd/ephemeral") {
    Write-Host "Installing from local source..." -ForegroundColor Blue
    go install ./cmd/ephemeral
} else {
    Write-Host "Installing from remote source..." -ForegroundColor Blue
    # Change the URL below to your actual repository URL
    go install github.com/yourusername/ephemeral/cmd/ephemeral@latest
}

# PATH check
$gopath = go env GOPATH
$gobin = Join-Path $gopath "bin"

if ($env:PATH -notlike "*$gobin*") {
    Write-Host "Warning: $gobin is not in your PATH." -ForegroundColor Red
    Write-Host "To add it temporarily: `$env:PATH += ';$gobin'"
    Write-Host "To add it permanently: [Environment]::SetEnvironmentVariable('Path', [Environment]::GetEnvironmentVariable('Path', 'User') + ';$gobin', 'User')"
}

Write-Host "Ephemeral installed successfully!" -ForegroundColor Green
Write-Host "Try running: ephemeral --nick Alice" -ForegroundColor Blue
