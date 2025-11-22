#!/bin/bash

set -e

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

packager() {
    sfx="${1:-}"
    filebrowser_dir="$filebrowser_os-$filebrowser_arch-filebrowser$sfx"
    filebrowser_file="$filebrowser_dir$filebrowser_dl_ext"

    echo "Packaging $filebrowser_bin for $filebrowser_dir"

    mkdir -p $filebrowser_dir
    cp $filebrowser_bin $filebrowser_dir/$filebrowser_bin
    tar -zcvf $filebrowser_file $filebrowser_dir

    if [ -f "$filebrowser_file" ]; then
        echo "Packaged $filebrowser_file successfully."
    else
        echo "Aborted, failed to package $filebrowser_file."
        return 8
    fi

    echo "Created archive: $filebrowser_file"
}

uploader() {
    echo "Uploading $filebrowser_file to GitHub Release"

    if [ -z "$GIT_TOKEN" ] || [ -z "$GITHUB_REPOSITORY" ]; then
        echo "GIT_TOKEN, or GITHUB_REPOSITORY is not set, skipping upload."
        return 0
    fi

    if [ -z "$RELEASE_ID" ]; then
        if [ -z "$TAG" ]; then
            echo "TAG is not set, attempting to use latest release."
            RELEASE_TAG="latest"
        else
            RELEASE_TAG="tags/$TAG"
            echo "TAG is set to '$TAG', attempting to get release ID for tag."
        fi
        RELEASE_ID=$(
            curl -s \
                -H "Accept: application/vnd.github+json" \
                -H "Authorization: Bearer $GIT_TOKEN" \
                https://api.github.com/repos/$GITHUB_REPOSITORY/releases/$RELEASE_TAG \
                | jq '.id'
        )
    fi

    if [ -z "$RELEASE_ID" ]; then
        echo "Aborted, RELEASE_ID is not set and could not be determined."
        return 9
    fi

    curl -X POST -H "Authorization: token $GIT_TOKEN" \
        -H "Content-Type: application/octet-stream" \
        --data-binary @"$filebrowser_file" \
        "https://uploads.github.com/repos/$GITHUB_REPOSITORY/releases/$RELEASE_ID/assets?name=$filebrowser_file"
}

packager $1
uploader
