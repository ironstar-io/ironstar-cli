#!/bin/bash

set -e
trap 'catch $? $LINENO' EXIT
catch() {
  if [ "$1" != "0" ]; then
    # error handling goes here
    echo "Error occurred during installation. Exiting..."
  fi
}

# Check the hardware type using uname
HARDWARE=$(uname -m)

# Check if the hardware is ARM or i386
if [ "$HARDWARE" == "arm64" ]; then
  # Running on an ARM processor"
  BINARY="iron-macos-arm64"
elif [ "$HARDWARE" == "x86_64" ]; then
  # Running on x86"
  BINARY="iron-macos"
else
  # Do something for other hardware types
  echo "Running on an unknown hardware type: $HARDWARE. Exiting..."
  exit 1
fi

# Set the URL of the GitHub releases page for the CLI
RELEASES_PAGE="https://github.com/ironstar-io/ironstar-cli/releases"

# Set the URL of the CLI binary in the latest release
BINARY_URL="https://github.com/ironstar-io/ironstar-cli/releases/latest/download/$BINARY"


install() {
    echo "Downloading the Ironstar CLI..."

    # Download the CLI binary and save it to the current directory
    curl -sL $BINARY_URL -o iron

    # Make the CLI binary executable
    chmod +x iron

    # Move the CLI binary to a directory in the PATH
    mv iron /usr/local/bin/iron

    echo "CLI installed successfully!"
}

install
