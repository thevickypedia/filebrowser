#!/bin/bash

set -e

rm -f filebrowser filebrowser.db filebrowser.exe

# read -p "Basic cleanup done. Continue with deep clean? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
echo -e "Basic cleanup done.\n"
read -p "Continue with deep clean? (Y/N): " confirm
[[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1

rm -rf frontend/node_modules frontend/dist
mkdir -p frontend/dist && touch frontend/dist/.gitkeep
