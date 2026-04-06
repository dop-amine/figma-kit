#!/bin/sh
set -e

REPO="amine/figma-kit"
BINARY="figma-kit"
INSTALL_DIR="${INSTALL_DIR:-/usr/local/bin}"

OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
ARCH="$(uname -m)"

case "$ARCH" in
  x86_64)  ARCH="amd64" ;;
  aarch64) ARCH="arm64" ;;
  arm64)   ARCH="arm64" ;;
  *)       echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

case "$OS" in
  darwin|linux) ;;
  *)            echo "Unsupported OS: $OS"; exit 1 ;;
esac

LATEST="$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')"
VERSION="${LATEST#v}"

ARCHIVE="${BINARY}_${VERSION}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$LATEST/$ARCHIVE"

echo "Downloading $BINARY v$VERSION for $OS/$ARCH..."
TMPDIR="$(mktemp -d)"
curl -fsSL "$URL" -o "$TMPDIR/$ARCHIVE"
tar -xzf "$TMPDIR/$ARCHIVE" -C "$TMPDIR"

echo "Installing to $INSTALL_DIR/$BINARY..."
sudo install -m 755 "$TMPDIR/$BINARY" "$INSTALL_DIR/$BINARY"
rm -rf "$TMPDIR"

echo "Done! Run 'figma-kit --version' to verify."
