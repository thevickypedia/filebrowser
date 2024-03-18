#!/usr/bin/env bash
trap 'echo -e "Aborted, error $? in command: $BASH_COMMAND"; trap ERR; return 1' ERR
filebrowser_os="unsupported"
filebrowser_dl_ext=".tar.gz"
filebrowser_bin="filebrowser"
unameu="$(tr '[:lower:]' '[:upper:]' <<<$(uname))"
if [[ $unameu == *DARWIN* ]]; then
  filebrowser_os="Darwin-x86_64"
elif [[ $unameu == *LINUX* ]]; then
  filebrowser_os="Linux-x86_64"
elif [[ $unameu == *WIN* || $unameu == MSYS* ]]; then
  # Should catch cygwin
  filebrowser_os="Windows-x86_64"
  filebrowser_bin+=".exe"
  filebrowser_dl_ext=".zip"
else
  echo "Aborted, unsupported or unknown OS: $unameu"
  return 6
fi
filebrowser_pkg="FileBrowser-${filebrowser_os}"
filebrowser_zipfile="$filebrowser_pkg$filebrowser_dl_ext"

mkdir -p "$filebrowser_pkg"
cp "$filebrowser_bin" "$filebrowser_pkg"/"$filebrowser_bin"
cp "README.md" "$filebrowser_pkg"/README.md
tar -zcvf "$filebrowser_zipfile" $filebrowser_pkg/
