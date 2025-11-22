#!/bin/bash

set -e

ORG=${ORG:-thevickypedia}
VERSION=${VERSION:-latest}
INSTALL_PATH=${INSTALL_PATH:-$(pwd)}

filebrowser_os="unsupported"
filebrowser_arch="unknown"

#########################
# Which OS and version? #
#########################

filebrowser_bin="filebrowser"
filebrowser_dl_ext=".tar.gz"

# NOTE: `uname -m` is more accurate and universal than `arch`
# See https://en.wikipedia.org/wiki/Uname
unamem="$(uname -m)"
case $unamem in
*aarch64*|arm64)
    filebrowser_arch="arm64";;
*64*)
    filebrowser_arch="amd64";;
*86*)
    filebrowser_arch="386";;
*armv5*)
    filebrowser_arch="armv5";;
*armv6*)
    filebrowser_arch="armv6";;
*armv7*)
    filebrowser_arch="armv7";;
*)
    echo "Aborted, unsupported or unknown architecture: $unamem"
    return 2
    ;;
esac

# shellcheck disable=SC2046
unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
if [[ $unameu == *DARWIN* ]]; then
    filebrowser_os="darwin"
elif [[ $unameu == *LINUX* ]]; then
    filebrowser_os="linux"
elif [[ $unameu == *FREEBSD* ]]; then
    filebrowser_os="freebsd"
elif [[ $unameu == *NETBSD* ]]; then
    filebrowser_os="netbsd"
elif [[ $unameu == *OPENBSD* ]]; then
    filebrowser_os="openbsd"
elif [[ $unameu == *WIN* || $unameu == MSYS* ]]; then
    # Should catch cygwin
    filebrowser_os="windows"
    filebrowser_bin="filebrowser.exe"
    filebrowser_dl_ext=".zip"
else
    echo "Aborted, unsupported or unknown OS: $(uname)"
    return 6
fi

########################
# Download and extract #
########################

if type -p curl >/dev/null 2>&1; then
    net_getter="curl -fsSL"
elif type -p wget >/dev/null 2>&1; then
    net_getter="wget -qO-"
else
    echo "Aborted, could not find curl or wget"
    return 7
fi

filebrowser_dir="$filebrowser_os-$filebrowser_arch-filebrowser"
filebrowser_file="$filebrowser_dir$filebrowser_dl_ext"
# shellcheck disable=SC2153
git_token="$GIT_TOKEN"
if [ -z "$git_token" ]; then
    headers=()
else
    headers=(-H "Authorization: Bearer $git_token")
fi

if [ "$VERSION" == "latest" ]; then
    release_url="https://api.github.com/repos/$ORG/filebrowser/releases/$VERSION"
else
    release_url="https://api.github.com/repos/$ORG/filebrowser/releases/tags/$VERSION"
fi
echo "Downloading filebrowser for $filebrowser_file from $release_url..."

response=$($net_getter "${headers[@]}" "$release_url")
asset_id=$(echo "$response" | jq -r --arg filebrowser_file "$filebrowser_file" '.assets[] | select(.name == $filebrowser_file) | .id')

if [ -z "$asset_id" ]; then
    echo "Failed to get the asset id for $filebrowser_file"
    exit 1
fi

asset_url="https://api.github.com/repos/$ORG/filebrowser/releases/assets/$asset_id"

$net_getter "${headers[@]}" -H "Accept: application/octet-stream" "$asset_url" > "$INSTALL_PATH/$filebrowser_file"

echo "Extracting..."
case "$filebrowser_file" in
    *.zip)    unzip -o "$filebrowser_file" ;;
    *.tar.gz) tar -xzf "$filebrowser_file" ;;
esac

echo "Moving 'filebrowser' to $INSTALL_PATH"
mv "$filebrowser_dir/$filebrowser_bin" "$INSTALL_PATH" && rm -rf $filebrowser_dir $filebrowser_file
