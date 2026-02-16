#!/bin/bash
# Ephemeral Installer for Linux, macOS, and Termux

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Starting Ephemeral Installation...${NC}"

# Check for Go
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed.${NC}"
    echo "Please install Go (1.23+) from https://golang.org/dl/"
    exit 1
fi

# Determine the source
# If we are in a git repo, use local. Otherwise, use remote.
if [ -d ".git" ] && [ -d "cmd/ephemeral" ]; then
    echo -e "${BLUE}Installing from local source...${NC}"
    go install ./cmd/ephemeral
else
    echo -e "${BLUE}Installing from remote source...${NC}"
    # Change the URL below to your actual repository URL
    go install github.com/yourusername/ephemeral/cmd/ephemeral@latest
fi

# Check if GOBIN is in PATH
GOBIN=$(go env GOPATH)/bin
if [[ ":$PATH:" != *":$GOBIN:"* ]]; then
    echo -e "${RED}Warning: $GOBIN is not in your PATH.${NC}"
    echo "Add the following line to your ~/.bashrc or ~/.zshrc:"
    echo "export PATH=\$PATH:$GOBIN"
fi

echo -e "${GREEN}Ephemeral installed successfully!${NC}"
echo -e "Try running: ${BLUE}ephemeral --nick Alice${NC}"
