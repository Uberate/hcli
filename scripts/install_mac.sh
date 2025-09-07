#!/bin/bash

# hcli macOS installation script
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
REPO_OWNER="uberate"  # Change this to your GitHub username
REPO_NAME="hugo-ai-helper"   # Change this to your repository name
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="/tmp/hcli-install"

# Detect architecture
ARCH=$(uname -m)
if [ "$ARCH" = "x86_64" ]; then
    ARCH="amd64"
elif [ "$ARCH" = "arm64" ]; then
    ARCH="arm64"
else
    echo -e "${RED}Unsupported architecture: $ARCH${NC}"
    exit 1
fi

# Detect macOS version
OS="darwin"

# Get latest release version
echo -e "${YELLOW}Fetching latest version...${NC}"
LATEST_VERSION=$(curl -s "https://api.github.com/repos/$REPO_OWNER/$REPO_NAME/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

if [ -z "$LATEST_VERSION" ]; then
    echo -e "${RED}Failed to fetch latest version${NC}"
    echo -e "${YELLOW}Please check:"
    echo -e "1. GitHub repository: https://github.com/$REPO_OWNER/$REPO_NAME"
    echo -e "2. Internet connection${NC}"
    exit 1
fi

echo -e "${GREEN}Latest version: $LATEST_VERSION${NC}"

# Download URL
DOWNLOAD_URL="https://github.com/$REPO_OWNER/$REPO_NAME/releases/download/$LATEST_VERSION/hcli_${OS}_${ARCH}.tar.gz"

# Create temporary directory
mkdir -p $TEMP_DIR
cd $TEMP_DIR

echo -e "${GREEN}Installing hcli v$VERSION for macOS ($ARCH)...${NC}"

# Download the release
echo -e "${YELLOW}Downloading hcli...${NC}"
if ! curl -L -o hcli.tar.gz "$DOWNLOAD_URL" --progress-bar; then
    echo -e "${RED}Failed to download hcli${NC}"
    echo -e "${YELLOW}Please check:"
    echo -e "1. The version number"
    echo -e "2. GitHub repository access"
    echo -e "3. Internet connection${NC}"
    exit 1
fi

# Extract the archive
echo -e "${YELLOW}Extracting...${NC}"
tar -xzf hcli.tar.gz

# Make binary executable
chmod +x hcli

# Install to /usr/local/bin
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if ! sudo mv hcli "$INSTALL_DIR/"; then
    echo -e "${RED}Failed to install to $INSTALL_DIR${NC}"
    echo -e "${YELLOW}You may need to run:${NC}"
    echo -e "  sudo mkdir -p $INSTALL_DIR"
    echo -e "  sudo chown $(whoami) $INSTALL_DIR"
    exit 1
fi

# Clean up
rm -rf $TEMP_DIR

# Verify installation
if command -v hcli &> /dev/null; then
    echo -e "${GREEN}Installation successful!${NC}"
    echo -e "${GREEN}hcli version:${NC}"
    hcli --version
else
    echo -e "${YELLOW}Installation completed but hcli not found in PATH${NC}"
    echo -e "${YELLOW}Please ensure $INSTALL_DIR is in your PATH${NC}"
fi

echo -e "${GREEN}Done! You can now use 'hcli' command.${NC}"