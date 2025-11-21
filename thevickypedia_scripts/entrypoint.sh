#!/bin/sh

# If config.json exists, import config
if [ -f "/config/config.json" ]; then
    echo "Importing config from /config/config.json"
    /filebrowser config import /config/config.json
fi

# If users.json exists, import users
if [ -f "/config/users.json" ]; then
    echo "Importing users from /config/users.json"
    /filebrowser users import /config/users.json
fi

# If the above settings are present, then the filebrowser.db might be at root
if [ -f "/filebrowser.db" ]; then
    echo "Database file exists at root, moving to /config."
    mv /filebrowser.db /config/filebrowser.db
fi

# Start the normal filebrowser server
exec /filebrowser \
    --root=/data \
    --database=/config/filebrowser.db
