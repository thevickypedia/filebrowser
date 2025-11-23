#!/bin/sh
set -e

# Config file
CONFIG_FILE=${CONFIG_FILE:-settings.json}   # default: settings.json
CONFIG_PATH="/config/$CONFIG_FILE"

if [ -f "$CONFIG_PATH" ]; then
    echo "Importing config from $CONFIG_PATH"
    /filebrowser config import "$CONFIG_PATH" >/dev/null 2>&1

    PORT=${FB_PORT:-$(jq -r '.server.port // empty' "$CONFIG_PATH")}
    ADDRESS=${FB_ADDRESS:-$(jq -r '.server.address // empty' "$CONFIG_PATH")}
fi

# Users file
USERS_FILE=${USERS_FILE:-users.json}       # default: users.json
USERS_PATH="/config/$USERS_FILE"

if [ -f "$USERS_PATH" ]; then
    echo "Importing users from $USERS_PATH"
    /filebrowser users import "$USERS_PATH" >/dev/null 2>&1
fi

# Move database if exists at root
if [ -f "/filebrowser.db" ]; then
    echo "Database file exists at root, moving to /config."
    mv /filebrowser.db /config/filebrowser.db
fi

# Defaults
export PORT=${PORT:-80}
export ADDRESS=${ADDRESS:-0.0.0.0}

echo "Starting Filebrowser with ADDRESS=$ADDRESS PORT=$PORT"
exec /filebrowser \
    --root=/data \
    --address="$ADDRESS" \
    --port="$PORT" \
    --database=/config/filebrowser.db
