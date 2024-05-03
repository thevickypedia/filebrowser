#!/bin/sh

if [ -n "$VOLUME_1" ]; then
    echo "Volume '$VOLUME_1' attached"
fi

/opt/filebrowser/filebrowser config import /opt/filebrowser/config.json
/opt/filebrowser/filebrowser users import /opt/filebrowser/users.json

/opt/filebrowser/filebrowser
