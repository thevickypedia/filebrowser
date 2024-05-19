import gzip
import logging
import os
import pathlib
import platform
import shutil
import stat
import tarfile
import zipfile

import requests

LOGGER = logging.getLogger(__name__)
HANDLER = logging.StreamHandler()
LOGGER.addHandler(HANDLER)
LOGGER.setLevel(logging.DEBUG)


class GitHub:
    """Custom GitHub account information loaded using multiple env prefixes.

    >>> GitHub

    """

    org: str = os.environ.get("ORG", "thevickypedia")
    token: str | None = os.environ.get("GIT_TOKEN")
    version: str = os.environ.get("git_version") or os.environ.get("GIT_VERSION") or "latest"


class Executable:
    """Executable object to load all the objects to download the executable from releases.

    >>> Executable

    """

    filebrowser_bin: str = "filebrowser"
    filebrowser_dl_ext: str = ".tar.gz"

    # Detect OS and architecture
    system: str = platform.system().lower()
    machine: str = platform.machine().lower()
    if system == "darwin":
        filebrowser_os: str = "darwin"
    elif system == "linux":
        filebrowser_os: str = "linux"
    elif system == "freebsd":
        filebrowser_os: str = "freebsd"
    elif system == "netbsd":
        filebrowser_os: str = "netbsd"
    elif system == "openbsd":
        filebrowser_os: str = "openbsd"
    elif system.startswith("win") or system == "msys":
        filebrowser_os: str = "windows"
        filebrowser_bin: str = "filebrowser.exe"
        filebrowser_dl_ext: str = ".zip"
    else:
        raise OSError(
            f"Aborted, unsupported or unknown OS: {system}"
        )

    if "aarch64" in machine or "arm64" in machine:
        filebrowser_arch: str = "arm64"
    elif "64" in machine:
        filebrowser_arch: str = "amd64"
    elif "86" in machine:
        filebrowser_arch: str = "386"
    elif "armv5" in machine:
        filebrowser_arch: str = "armv5"
    elif "armv6" in machine:
        filebrowser_arch: str = "armv6"
    elif "armv7" in machine:
        filebrowser_arch: str = "armv7"
    else:
        raise OSError(
            f"Aborted, unsupported or unknown architecture: {machine}"
        )

    filebrowser_file: str = f"{filebrowser_os}-{filebrowser_arch}-filebrowser{filebrowser_dl_ext}"
    filebrowser_db: str = f"{filebrowser_bin}.db"


github = GitHub()
executable = Executable()


def download():
    LOGGER.info("Source Repository: 'https://github.com/%s/filebrowser'", github.org)
    LOGGER.info("Targeted Asset: '%s'", executable.filebrowser_file)
    headers = {"Authorization": f"Bearer {github.token}"} if github.token else {}
    # Get the release from the specified version
    if github.version == "latest":
        release_url = f"https://api.github.com/repos/{github.org}/filebrowser/releases/{github.version}"
    else:
        release_url = f"https://api.github.com/repos/{github.org}/filebrowser/releases/tags/{github.version}"
    response = requests.get(release_url, headers=headers)
    response.raise_for_status()
    release_info = response.json()

    # Log the download URL
    filebrowser_url = f"https://github.com/{github.org}/filebrowser/releases/download/" \
                      f"{release_info['tag_name']}/{executable.filebrowser_file}"
    LOGGER.info("Download URL: %s", filebrowser_url)

    # Get asset id
    existing = []
    for asset in release_info['assets']:
        if asset.get('name') == executable.filebrowser_file:
            asset_id = asset['id']
            break
        elif asset.get('name'):
            existing.append(asset['name'])
    else:
        existing = '\n\t'.join(existing)
        raise Exception(
            f"\n\tFailed to get the asset id for {executable.filebrowser_file!r}\n\n"
            f"Available asset names:\n\t{existing}"
        )

    # Download the asset
    headers['Accept'] = "application/octet-stream"
    response = requests.get(f"https://api.github.com/repos/{github.org}/filebrowser/releases/assets/{asset_id}",
                            headers=headers)
    response.raise_for_status()
    with open(executable.filebrowser_file, 'wb') as file:
        for chunk in response.iter_content(chunk_size=8192):
            file.write(chunk)
    assert os.path.isfile(executable.filebrowser_file), f"Failed to get the asset id for {executable.filebrowser_file}"
    LOGGER.info("Asset has been downloaded successfully")

    # Extract asset based on the file extension
    if executable.filebrowser_file.endswith(".tar.gz"):
        tar_file = executable.filebrowser_file.removesuffix('.gz')
        # Read the gzipped file as bytes, and write as a tar file (.tar.gz -> .tar)
        with gzip.open(executable.filebrowser_file, 'rb') as f_in:
            with open(tar_file, 'wb') as f_out:
                f_out.write(f_in.read())
                f_out.flush()
        assert os.path.isfile(tar_file), f"Failed to gunzip {executable.filebrowser_file}"
        os.remove(executable.filebrowser_file)

        # Read the tar file and extract its content
        with tarfile.open(tar_file, 'r') as tar:
            tar.extractall()
        # Catches the use case where binary might be directly archived
        if os.path.isfile(executable.filebrowser_bin):
            return True
        content_dir = tar_file.removesuffix('.tar')
        assert os.path.isdir(content_dir) and os.path.isfile(os.path.join(content_dir, executable.filebrowser_bin)), \
            f"Failed to unarchive {tar_file}"
        os.remove(tar_file)
    elif executable.filebrowser_file.endswith(".zip"):
        # Read the zip file and extract its content
        with zipfile.ZipFile(executable.filebrowser_file, 'r') as zip_ref:
            zip_ref.extractall()
        # Catches the use case where binary might be directly zipped
        if os.path.isfile(executable.filebrowser_bin):
            return True
        content_dir = executable.filebrowser_file.removesuffix(".zip")
        assert os.path.isdir(content_dir) and os.path.isfile(os.path.join(content_dir, executable.filebrowser_bin)), \
            f"Failed to unzip {executable.filebrowser_file}"
        os.remove(executable.filebrowser_file)
    else:
        raise OSError(
            f"Invalid filename: {executable.filebrowser_file}"
        )

    # Copy the executable out of the extracted directory and remove the extraction directory
    shutil.copyfile(os.path.join(content_dir, executable.filebrowser_bin),
                    os.path.join(os.getcwd(), executable.filebrowser_bin))
    shutil.rmtree(content_dir)


if __name__ == '__main__':
    download()
    # Change file permissions and set as executable
    # os.chmod(executable.filebrowser_bin, 0o755)
    # basically, chmod +x => -rwxr-xr-x
    os.chmod(executable.filebrowser_bin, stat.S_IRWXU | stat.S_IRGRP | stat.S_IXGRP | stat.S_IROTH | stat.S_IXOTH)
    shutil.move(executable.filebrowser_bin, pathlib.Path(os.getcwd()).parent)
