#!/usr/bin/env bash

set -e

# Set parent directory as current working directory
CURRENT_PATH="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BASE_PATH="$(dirname $CURRENT_PATH)"
cd "$BASE_PATH"

make lint-frontend
make lint-backend
make test-frontend
make test-backend
