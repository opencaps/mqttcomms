#!/bin/bash
# Generate AsyncAPI HTML documentation
# Requires: npm install -g @asyncapi/cli

set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_ROOT"

echo "Generating AsyncAPI documentation..."

# Check if asyncapi CLI is installed
if ! command -v asyncapi &> /dev/null; then
    echo "AsyncAPI CLI not found. Installing..."
    npm install -g @asyncapi/cli
fi

# Generate HTML documentation
asyncapi generate fromTemplate asyncapi.yaml @asyncapi/html-template -o docs --force-write

echo "Documentation generated in docs/ directory"
echo "Open docs/index.html in your browser to view"
