#!/usr/bin/env bash

set -e

printer() {
  echo "************************************************************************************************************************************************"
  echo "$1"
  echo "************************************************************************************************************************************************"
  echo ""
}

cleanup() {
  rm -f filebrowser filebrowser.db filebrowser.exe
  rm -rf frontend/node_modules frontend/dist
  mkdir -p frontend/dist && touch frontend/dist/.gitkeep
}

# Set parent directory as current working directory
CURRENT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(dirname $CURRENT_PATH)"
cd "$BASE_PATH"

cleanup

SHELL="/usr/bin/env bash"
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

printer "Building filebrowser frontend..."
cd frontend && pnpm install --frozen-lockfile && pnpm run build

# Run on a new shell to avoid segmentation error
printer "Building filebrowser backend..."
bash -c "cd ${BASE_PATH} && go build -ldflags \"$LDFLAGS\" -o ."

printer "Completed build..."
