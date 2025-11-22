import sys
import os
import requests
import pathlib
import pprint

try:
    FROM_VERSION = sys.argv[1]
except IndexError:
    raise ValueError("Please provide the 'from' version as the first argument.")

try:
    TO_VERSION = sys.argv[2]
except IndexError:
    TO_VERSION = "master"

ROOT_DIR = pathlib.Path(__file__).parent.parent.resolve()

REPO_OWNER = "filebrowser"
REPO_NAME = "filebrowser"

COMPARE_URL = f"https://api.github.com/repos/{REPO_OWNER}/{REPO_NAME}/compare/{FROM_VERSION}...{TO_VERSION}"

print(f"Fetching compare data from:\n{COMPARE_URL}\n")
response = requests.get(COMPARE_URL)
data = response.json()

if "files" not in data:
    raise RuntimeError("Unexpected API response (missing 'files')")

# -----------------------------
# Categorize files
# -----------------------------
changes = {
    "changed": [],
    "added": [],
    "removed": [],
    "renamed": []
}

for f in data["files"]:
    status = f["status"]

    if status == "modified":
        changes["changed"].append(f)
    elif status == "added":
        changes["added"].append(f)
    elif status == "removed":
        changes["removed"].append(f)
    elif status == "renamed":
        changes["renamed"].append(f)

print("\nDetected changes summary:")
print(f"  Added:   {len(changes['added'])}")
print(f"  Changed: {len(changes['changed'])}")
print(f"  Removed: {len(changes['removed'])}")
print(f"  Renamed: {len(changes['renamed'])}")
print()
# -----------------------------

def download_file(file_info, base_path):
    """Download a file at specific commit into a path."""
    file_path = file_info["filename"]
    raw_url = file_info["raw_url"]

    out_path = os.path.join(base_path, file_path)
    os.makedirs(os.path.dirname(out_path), exist_ok=True)

    print(f"Downloading {file_path} → {out_path}")
    r = requests.get(raw_url)
    r.raise_for_status()

    with open(out_path, "wb") as f:
        f.write(r.content)
        f.flush()


# ------------------------------
# Download added + changed files
# ------------------------------
print("Downloading ADDED and CHANGED files...\n")

for file_info in changes["added"]:
    download_file(file_info, ROOT_DIR)

for file_info in changes["changed"]:
    download_file(file_info, ROOT_DIR)

print("\nDownload of added+changed files completed.\n")


# ---------------------------
# Auto-rename RENAMED files
# ---------------------------
if changes["renamed"]:
    print("\nAuto-renaming files...")

    for f in changes["renamed"]:
        old_path = os.path.join(ROOT_DIR, f["previous_filename"])
        new_path = os.path.join(ROOT_DIR, f["filename"])

        os.makedirs(os.path.dirname(new_path), exist_ok=True)

        if os.path.exists(old_path):
            print(f"Renaming: {old_path} → {new_path}")
            os.rename(old_path, new_path)
        else:
            print(f"[WARN] Old file missing, downloading instead: {f['previous_filename']}")
            download_file(f, ROOT_DIR)


# ------------------------
# Auto-remove REMOVED files
# ------------------------
if changes["removed"]:
    print("\nAuto-removing files...")

    for f in changes["removed"]:
        file_path = os.path.join(ROOT_DIR, f["filename"])

        if os.path.exists(file_path):
            print(f"Removing: {file_path}")
            os.remove(file_path)
        else:
            print(f"[WARN] File already missing: {file_path}")


print("\n=============================")
print("All operations complete!")
print(f"Files saved under: {ROOT_DIR}")
print("=============================\n")
