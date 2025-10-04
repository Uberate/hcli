#!/bin/bash

# Build and install script from source
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuration
INSTALL_DIR="/usr/local/bin"
TEMP_DIR="./tmp/hcli-build"

# Parse OS type
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
case $OS in
    "linux")
        OS="linux"
        ;;
    "darwin")
        OS="darwin"
        ;;
    *)
        echo -e "${RED}Unsupported operating system: $OS${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Detected OS: $OS${NC}"

# Parse CPU architecture
ARCH=$(uname -m)
case $ARCH in
    "x86_64")
        ARCH="amd64"
        ;;
    "amd64")
        ARCH="amd64"
        ;;
    "arm64" | "aarch64")
        ARCH="arm64"
        ;;
    "armv7l")
        ARCH="armv7"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

echo -e "${GREEN}Detected architecture: $ARCH${NC}"

# Build binary using make releases
echo -e "${YELLOW}Building binary with make releases...${NC}"
if ! make releases; then
    echo -e "${RED}Build failed${NC}"
    exit 1
fi

# Find the built binary
BUILD_OUTPUT="output/hcli_${OS}_${ARCH}"
if [ ! -d "$BUILD_OUTPUT" ]; then
    echo -e "${RED}Build output directory not found: $BUILD_OUTPUT${NC}"
    echo -e "${YELLOW}Available build outputs:${NC}"
    ls -la output/ 2>/dev/null || echo "No output directory found"
    exit 1
fi

BINARY_PATH="$BUILD_OUTPUT/hcli"
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Binary not found at: $BINARY_PATH${NC}"
    echo -e "${YELLOW}Available files:${NC}"
    ls -la "$BUILD_OUTPUT/"
    exit 1
fi

# Make binary executable
chmod +x "$BINARY_PATH"

# Install to /usr/local/bin
echo -e "${YELLOW}Installing to $INSTALL_DIR...${NC}"
if ! sudo mv "$BINARY_PATH" "$INSTALL_DIR/"; then
    echo -e "${RED}Failed to install to $INSTALL_DIR${NC}"
    echo -e "${YELLOW}You may need to:${NC}"
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

echo -e ""
echo -e "${GREEN}Done! hcli has been built from source and installed to $INSTALL_DIR${NC}"