#!/bin/bash

set -e

rm -f filebrowser filebrowser.db auth.db filebrowser.exe

# read -p "Basic cleanup done. Continue with deep clean? (Y/N): " confirm && [[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1
echo -e "Basic cleanup done.\n"
# check if the script was triggered with -y or --yes flag
if [[ "$1" == "-y" || "$1" == "--yes" ]]; then
    confirm="y"
else
    read -p "Continue with deep clean? (Y/N): " confirm
fi
[[ $confirm == [yY] || $confirm == [yY][eE][sS] ]] || exit 1

rm -rf frontend/node_modules frontend/dist vendor/
mkdir -p frontend/dist && touch frontend/dist/.gitkeep

echo -e "Deep cleanup done.\n"
