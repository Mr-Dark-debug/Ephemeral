#!/bin/bash
# Ephemeral Installer for Linux, macOS, and Termux

set -e

# Configuration
REPO_OWNER="Mr-Dark-debug"
REPO_NAME="Ephemeral"
BINARY_NAME="ephemeral"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

echo -e "${BLUE}Starting Ephemeral Installation...${NC}"

# Detect OS and Arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo -e "${RED}Unsupported architecture: $ARCH${NC}"; exit 1 ;;
esac

# Special case for android
if [ "$OS" = "linux" ] && [ -d "/data/data/com.termux" ]; then
    OS="android"
fi

echo -e "Detected: ${GREEN}${OS}/${ARCH}${NC}"

# Fetch latest release URL
API_URL="https://api.github.com/repos/${REPO_OWNER}/${REPO_NAME}/releases/latest"
DOWNLOAD_URL=$(curl -s $API_URL | grep "browser_download_url" | grep "${OS}_${ARCH}" | cut -d '"' -f 4 | head -n 1)

if [ -z "$DOWNLOAD_URL" ]; then
    echo -e "${RED}Error: Could not find a release for your platform.${NC}"
    exit 1
fi

# Download and Install
INSTALL_DIR="${HOME}/.local/bin"
mkdir -p "$INSTALL_DIR"

TEMP_FILE=$(mktemp)
echo -e "Downloading from ${BLUE}${DOWNLOAD_URL}${NC}..."
curl -L -o "$TEMP_FILE" "$DOWNLOAD_URL"

echo -e "Installing to ${GREEN}${INSTALL_DIR}/${BINARY_NAME}${NC}..."

# Extract based on extension
if [[ "$DOWNLOAD_URL" == *.tar.gz ]]; then
    tar -xzf "$TEMP_FILE" -C "$INSTALL_DIR" "$BINARY_NAME"
elif [[ "$DOWNLOAD_URL" == *.zip ]]; then
    # Termux might not have unzip by default
    if ! command -v unzip &> /dev/null; then
        echo "Installing unzip..."
        if [ "$OS" = "android" ]; then pkg install unzip -y; fi
    fi
    unzip -o "$TEMP_FILE" "$BINARY_NAME" -d "$INSTALL_DIR"
else
    # Direct binary
    mv "$TEMP_FILE" "$INSTALL_DIR/${BINARY_NAME}"
fi

chmod +x "$INSTALL_DIR/${BINARY_NAME}"
rm -f "$TEMP_FILE"

# PATH check
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${RED}Warning: $INSTALL_DIR is not in your PATH.${NC}"
    echo "Add the following to your ~/.bashrc or ~/.zshrc:"
    echo "export PATH=\$PATH:$INSTALL_DIR"
fi

echo -e "${GREEN}Ephemeral installed successfully!${NC}"
echo -e "Try running: ${BLUE}ephemeral --nick Alice${NC}"
