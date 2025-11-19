#!/bin/bash
# Install man page for pk

set -e

MAN_FILE="docs/pk.1"
MAN_DIR="/usr/local/share/man/man1"

if [ ! -f "$MAN_FILE" ]; then
    echo "Error: $MAN_FILE not found"
    exit 1
fi

echo "Installing pk man page..."

# Create man directory if it doesn't exist
sudo mkdir -p "$MAN_DIR"

# Copy man page
sudo cp "$MAN_FILE" "$MAN_DIR/pk.1"

# Update man database
if command -v mandb &> /dev/null; then
    sudo mandb -q
fi

echo "✓ Man page installed to $MAN_DIR/pk.1"
echo ""
echo "View with:"
echo "  man pk"
