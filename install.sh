#!/bin/sh
set -e

REPO="KristianJBorgwarth/dendrite.daemon"
BINARY="dendrite"

# Detect OS
OS="$(uname -s)"
case "$OS" in
  Linux)  GOOS="linux" ;;
  Darwin) GOOS="darwin" ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

# Detect architecture
ARCHITECTURE="$(uname -m)"
case "$ARCHITECTURE" in
  x86_64)          GOARCH="amd64" ;;
  arm64 | aarch64) GOARCH="arm64" ;;
  *)
    echo "Unsupported architecture: $ARCHITECTURE"
    exit 1
    ;;
esac

# Resolve install dir
if [ "$GOOS" = "darwin" ]; then
  INSTALL_DIR="/usr/local/bin"
else
  INSTALL_DIR="$HOME/.local/bin"
  mkdir -p "$INSTALL_DIR"
fi

# Fetch latest release tag
echo "Fetching latest release..."
TAG="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" \
  | grep '"tag_name"' \
  | sed 's/.*"tag_name": *"\(.*\)".*/\1/')"

if [ -z "$TAG" ]; then
  echo "Could not determine latest release tag."
  exit 1
fi

ASSET="${BINARY}-${GOOS}-${GOARCH}"
URL="https://github.com/$REPO/releases/download/$TAG/$ASSET"

echo "Downloading $ASSET ($TAG)..."
TMP="$(mktemp)"
curl -fsSL "$URL" -o "$TMP"
chmod +x "$TMP"
mv "$TMP" "$INSTALL_DIR/$BINARY"

echo ""
echo "dendrite $TAG installed to $INSTALL_DIR/$BINARY"

# Warn if install dir is not on PATH
case ":$PATH:" in
  *":$INSTALL_DIR:"*) ;;
  *)
    echo ""
    echo "Warning: $INSTALL_DIR is not on your PATH."
    echo "Add the following to your shell profile:"
    echo ""
    echo "  export PATH=\"$INSTALL_DIR:\$PATH\""
    ;;
esac
