#!/bin/sh

# If config.json exists, import config
if [ -f "/config/settings.json" ]; then
    echo "Importing config from /config/settings.json"
    /filebrowser config import /config/settings.json >/dev/null 2>&1
    PORT=${FB_PORT:-$(jq -r .port /config/settings.json)}
    ADDRESS=${FB_ADDRESS:-$(jq -r .address /config/settings.json)}
fi

# If users.json exists, import users
if [ -f "/config/users.json" ]; then
    echo "Importing users from /config/users.json"
    /filebrowser users import /config/users.json >/dev/null 2>&1
fi

# If the above settings are present, then the filebrowser.db might be at root
if [ -f "/filebrowser.db" ]; then
    echo "Database file exists at root, moving to /config."
    mv /filebrowser.db /config/filebrowser.db
fi

export PORT=${PORT:-80}
export ADDRESS=${ADDRESS:-0.0.0.0}

# Start the normal filebrowser server
exec /filebrowser \
    --root=/data \
    --address=$ADDRESS \
    --port=$PORT \
    --database=/config/filebrowser.db
