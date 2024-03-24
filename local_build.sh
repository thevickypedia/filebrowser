#!/usr/bin/env bash

set -e

cleanup() {
  rm -rf filebrowser filebrowser.db frontend/node_modules
  rm -rf frontend/dist && mkdir -p frontend/dist && touch frontend/dist/.gitkeep
}

cleanup &

SHELL="/usr/bin/env bash"
BASE_PATH="$(dirname "$(realpath "${BASH_SOURCE[0]}")")"
VERSION=$(git describe --tags --always --match=v* 2> /dev/null || cat "$BASE_PATH/.version" 2> /dev/null || echo v0)
VERSION_HASH=$(git rev-parse HEAD)

go() {
    GOGC=off go "$@"
}

MODULE=$(env GO111MODULE=on go list -m)
TOOLS_DIR="$BASE_PATH/tools"
TOOLS_BIN="$TOOLS_DIR/bin"
mkdir -p "$TOOLS_BIN"
export PATH="$TOOLS_BIN:$PATH"
LDFLAGS+="-X \"$MODULE/version.Version=$VERSION\" -X \"$MODULE/version.CommitSHA=$VERSION_HASH\""

CURRENT_DIR="$(pwd)"

cd frontend && npm ci && npm run build

# Run on a new shell to avoid segmentation error
bash -c "cd $CURRENT_DIR && go build -ldflags \"$LDFLAGS\" -o ."
