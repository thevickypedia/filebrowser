#!/bin/bash

# Get the absolute path of the git repository root
GIT_ROOT=$(git rev-parse --show-toplevel)

# Construct the target path relative to the root
TARGET_DIR="$GIT_ROOT/frontend/src/i18n"

for file in $(git diff --name-only -- "$TARGET_DIR"/*.json); do
  old_value=$(git show HEAD:"$file" 2>/dev/null | jq -r '.login.otpPlaceholder' 2>/dev/null)

  if [ -n "$old_value" ] && [ "$old_value" != "null" ]; then
    jq --arg val "$old_value" '
      .login = (.login // {}) |
      .login.otpPlaceholder = $val
    ' "$file" > tmp.$$.json && mv tmp.$$.json "$file"
  fi
done
