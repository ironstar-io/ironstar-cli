#!/bin/bash

set -e
trap 'catch $? $LINENO' EXIT
catch() {
  if [ "$1" != "0" ]; then
    echo "Error occurred during installation. Exiting..."
  fi
}

# The macOS release is a single signed, notarised universal binary (Intel +
# Apple Silicon), so one asset serves every Mac.
BINARY="iron-macos"

# Set the URL of the CLI binary in the latest release
BINARY_URL="https://github.com/ironstar-io/ironstar-cli/releases/latest/download/$BINARY"

install() {
    echo "Downloading the Ironstar CLI..."

    # Download the CLI binary and save it to the current directory
    curl -sL "$BINARY_URL" -o iron

    # Make the CLI binary executable
    chmod +x iron

    # Move the CLI binary to a directory in the PATH
    mv iron /usr/local/bin/iron

    echo "CLI installed successfully!"
}

install
