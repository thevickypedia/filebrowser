#!/bin/sh

set -e

# Config file
CONFIG_FILE=${CONFIG_FILE:-settings.json}   # default: settings.json
CONFIG_PATH="/config/$CONFIG_FILE"

if [ -f "$CONFIG_PATH" ]; then
    PORT=${FB_PORT:-$(jq -r '.server.port // empty' "$CONFIG_PATH")}
    ADDRESS=${FB_ADDRESS:-$(jq -r '.server.address // empty' "$CONFIG_PATH")}
fi

# Defaults
export PORT=${PORT:-80}
export ADDRESS=${ADDRESS:-0.0.0.0}

wget -q --spider http://$ADDRESS:$PORT/health || exit 1
