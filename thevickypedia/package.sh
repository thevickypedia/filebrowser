#!/usr/bin/env bash

set -e

# Set parent directory as current working directory
parent="$(dirname "$PWD")"
cd "$parent"

package_filebrowser()
{
  trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; return 1' ERR
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
  *aarch64*)
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

  if [ -e "$filebrowser_bin" ]; then
      echo "Packaging executable '$filebrowser_bin'"
  else
      echo ""
      echo "***************************************************************************************************************";
      echo "                                  Executable '$filebrowser_bin' doesn't exist                                  ";
      echo "                                 Please run 'build.sh' to create the executable                                ";
      echo "***************************************************************************************************************";
      echo ""
      exit 1
  fi

  filebrowser_pkg="$filebrowser_os-$filebrowser_arch-filebrowser"
  filebrowser_zipfile="$filebrowser_os-$filebrowser_arch-filebrowser$filebrowser_dl_ext"

  mkdir -p "$filebrowser_pkg"
  cp "$filebrowser_bin" "$filebrowser_pkg/$filebrowser_bin"
  cp "README.md" "$filebrowser_pkg"/README.md
  tar -zcvf "$filebrowser_zipfile" $filebrowser_pkg/
  rm -rf "$filebrowser_pkg"
}

package_filebrowser
