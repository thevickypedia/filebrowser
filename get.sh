#!/usr/bin/env bash
trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; return 1' ERR
filemanager_os="unsupported"
filemanager_arch="unknown"
install_path="/usr/local/bin"

# Termux on Android has $PREFIX set which already ends with /usr
if [[ -n "$ANDROID_ROOT" && -n "$PREFIX" ]]; then
  install_path="$PREFIX/bin"
fi

# Fall back to /usr/bin if necessary
if [[ ! -d $install_path ]]; then
  install_path="/usr/bin"
fi

# Not every platform has or needs sudo (https://termux.com/linux.html)
((EUID)) && [[ -z "$ANDROID_ROOT" ]] && sudo_cmd="sudo"

#########################
# Which OS and version? #
#########################

filemanager_bin="filebrowser"
filemanager_dl_ext=".tar.gz"

# NOTE: `uname -m` is more accurate and universal than `arch`
# See https://en.wikipedia.org/wiki/Uname
unamem="$(uname -m)"
case $unamem in
*aarch64*)
  filemanager_arch="arm64";;
*64*)
  filemanager_arch="amd64";;
*86*)
  filemanager_arch="386";;
*armv5*)
  filemanager_arch="armv5";;
*armv6*)
  filemanager_arch="armv6";;
*armv7*)
  filemanager_arch="armv7";;
*)
  echo "Aborted, unsupported or unknown architecture: $unamem"
  return 2
  ;;
esac

unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
if [[ $unameu == *DARWIN* ]]; then
  filemanager_os="darwin"
elif [[ $unameu == *LINUX* ]]; then
  filemanager_os="linux"
elif [[ $unameu == *FREEBSD* ]]; then
  filemanager_os="freebsd"
elif [[ $unameu == *NETBSD* ]]; then
  filemanager_os="netbsd"
elif [[ $unameu == *OPENBSD* ]]; then
  filemanager_os="openbsd"
elif [[ $unameu == *WIN* || $unameu == MSYS* ]]; then
  # Should catch cygwin
  sudo_cmd=""
  filemanager_os="windows"
  filemanager_bin="filebrowser.exe"
  filemanager_dl_ext=".zip"
else
  echo "Aborted, unsupported or unknown OS: $uname"
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

filemanager_file="${filemanager_os}-$filemanager_arch-filebrowser$filemanager_dl_ext"
filemanager_tag="$(${net_getter} -H "Authorization: Bearer $GIT_TOKEN" https://api.github.com/repos/thevickypedia/filebrowser/releases/latest | grep -o '"tag_name": ".*"' | sed 's/"//g' | sed 's/tag_name: //g')"
filemanager_url="https://github.com/thevickypedia/filebrowser/releases/download/$filemanager_tag/$filemanager_file"
echo "$filemanager_url"
