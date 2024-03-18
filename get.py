import gzip
import logging
import os
import platform
import shutil
import stat
import tarfile
import zipfile

import requests

EXECUTABLE = "filebrowser"

logger = logging.getLogger(__name__)
logger.setLevel(level=logging.DEBUG)
handler = logging.StreamHandler()
handler.setFormatter(
    fmt=logging.Formatter(
        fmt='%(asctime)s - %(levelname)s - [%(processName)s:%(module)s:%(lineno)d] - %(funcName)s - %(message)s'
    )
)
logger.addHandler(hdlr=handler)


def download_asset() -> None:
    """Downloads the latest released asset."""
    global EXECUTABLE

    filemanager_dl_ext = ".tar.gz"

    # Detect OS
    system = platform.system().lower()
    if system == "darwin":
        filemanager_os = "Darwin-x86_64"
    elif system == "linux":
        filemanager_os = "Linux-x86_64"
    elif system.startswith("win") or system == "msys":
        filemanager_os = "Windows-x86_64"
        filemanager_dl_ext = ".zip"
    else:
        raise Exception(f"Aborted, unsupported or unknown OS: {system}")

    filemanager_file = f"FileBrowser-{filemanager_os}{filemanager_dl_ext}"
    git_token = os.environ.get("GIT_TOKEN")
    headers = {"Authorization": f"Bearer {git_token}"} if git_token else {}
    response = requests.get("https://api.github.com/repos/thevickypedia/filebrowser/releases/latest",
                            headers=headers)
    response.raise_for_status()

    release_info = response.json()
    for asset in release_info['assets']:
        if asset.get('name') == filemanager_file:
            asset_id = asset['id']
            break
    else:
        raise Exception(f"Failed to get the asset id for {filemanager_file}")

    headers['Accept'] = "application/octet-stream"
    response = requests.get(f"https://api.github.com/repos/thevickypedia/filebrowser/releases/assets/{asset_id}",
                            headers=headers)
    response.raise_for_status()

    with open(filemanager_file, 'wb') as file:
        for chunk in response.iter_content(chunk_size=8192):
            file.write(chunk)
    assert os.path.isfile(filemanager_file), f"Failed to get the asset id for {filemanager_file}"
    logger.info("Asset has been downloaded successfully")

    if filemanager_file.endswith(".tar.gz"):
        tar_file = filemanager_file.rstrip('.gz')
        with gzip.open(filemanager_file, 'rb') as f_in:
            with open(tar_file, 'wb') as f_out:
                f_out.write(f_in.read())
                f_out.flush()
        assert os.path.isfile(tar_file), f"Failed to gunzip {filemanager_file}"
        os.remove(filemanager_file)
        content_dir = tar_file.rstrip('.tar')
        with tarfile.open(tar_file, 'r') as tar:
            tar.extractall()
        assert os.path.isdir(content_dir) and os.path.isfile(os.path.join(content_dir, EXECUTABLE)), \
            f"Failed to unarchive {tar_file}"
        os.remove(tar_file)
    elif filemanager_file.endswith(".zip"):
        EXECUTABLE += ".exe"
        content_dir = filemanager_file.rstrip(".zip")
        with zipfile.ZipFile(filemanager_file, 'r') as zip_ref:
            zip_ref.extractall()
        assert os.path.isdir(content_dir) and os.path.isfile(os.path.join(content_dir, EXECUTABLE)), \
            f"Failed to unzip {filemanager_file}"
        os.remove(filemanager_file)
    else:
        raise OSError(
            f"Invalid filename: {filemanager_file}"
        )

    shutil.copyfile(os.path.join(content_dir, EXECUTABLE), os.path.join(os.getcwd(), EXECUTABLE))
    shutil.rmtree(content_dir)

    # os.chmod(EXECUTABLE, 0o755)
    # basically, chmod +x => -rwxr-xr-x
    os.chmod(EXECUTABLE, stat.S_IRWXU | stat.S_IRGRP | stat.S_IXGRP | stat.S_IROTH | stat.S_IXOTH)

    logger.info(f"Asset {EXECUTABLE!r} is ready to be used")


if __name__ == '__main__':
    download_asset()
