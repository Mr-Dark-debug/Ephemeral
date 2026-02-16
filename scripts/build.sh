#!/bin/bash
set -e
mkdir -p dist
echo "Building for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o dist/ephemeral-linux-amd64 ./cmd/ephemeral
echo "Building for Linux (arm64)..."
GOOS=linux GOARCH=arm64 go build -o dist/ephemeral-linux-arm64 ./cmd/ephemeral
echo "Building for macOS (amd64)..."
GOOS=darwin GOARCH=amd64 go build -o dist/ephemeral-darwin-amd64 ./cmd/ephemeral
echo "Building for Windows (amd64)..."
GOOS=windows GOARCH=amd64 go build -o dist/ephemeral-windows-amd64.exe ./cmd/ephemeral
echo "Building for Android (arm64)..."
GOOS=android GOARCH=arm64 go build -o dist/ephemeral-android-arm64 ./cmd/ephemeral
echo "Build complete. Artifacts in dist/"
