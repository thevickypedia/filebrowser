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
REMOVE_FILES = (
    "CHANGELOG.md",
    "CODE-OF-CONDUCT.md",
    "CONTRIBUTING.md",
    "SECURITY.md",
    "Taskfile.yml",
    "transifex.yml",
    ".goreleaser.yml",
    "renovate.json",
    ".github/CODEOWNERS",
    ".github/PULL_REQUEST_TEMPLATE.md",
    ".github/ISSUE_TEMPLATE/bug_report.yml",
    ".github/ISSUE_TEMPLATE/config.yml",
    ".github/workflows/ci.yaml",
    ".github/workflows/docs.yml",
)

IGNORE_FILES = (
    "README.md",
)

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


def should_ignore(file_path: str) -> bool:
    """Check if a file should be ignored based on its path."""
    return file_path in IGNORE_FILES


def rewrite_imports(path: str):
    """Replace github.com/filebrowser with github.com/thevickypedia in a file."""
    try:
        with open(path, "r", encoding="utf-8") as f:
            content = f.read()
    except UnicodeDecodeError:
        # Skip binary files (images etc.)
        return

    new_content = content.replace("github.com/filebrowser", "github.com/thevickypedia")

    if new_content != content:
        print(f"Rewriting imports in: {path}")
        with open(path, "w", encoding="utf-8") as f:
            f.write(new_content)


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

    # Rewrite import paths if it's a text file
    rewrite_imports(out_path)


# ------------------------------
# Download added + changed files
# ------------------------------
print("Downloading ADDED and CHANGED files...\n")

for file_info in changes["added"]:
    if should_ignore(file_info["filename"]):
        print(f"Ignoring (added): {file_info['filename']}")
        continue
    download_file(file_info, ROOT_DIR)

for file_info in changes["changed"]:
    if should_ignore(file_info["filename"]):
        print(f"Ignoring (changed): {file_info['filename']}")
        continue
    download_file(file_info, ROOT_DIR)

print("\nDownload of added+changed files completed.\n")


# ---------------------------
# Auto-rename RENAMED files
# ---------------------------
if changes["renamed"]:
    print("\nAuto-renaming files...")

    for f in changes["renamed"]:
        if should_ignore(f["filename"]) or should_ignore(f["previous_filename"]):
            print(f"Ignoring (renamed): {f['previous_filename']} → {f['filename']}")
            continue

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
        if should_ignore(f["filename"]):
            print(f"Ignoring (removed): {f['filename']}")
            continue
        file_path = os.path.join(ROOT_DIR, f["filename"])

        if os.path.exists(file_path):
            print(f"Removing: {file_path}")
            os.remove(file_path)
        else:
            print(f"[WARN] File already missing: {file_path}")

# ------------------------
# Remove unnecessary files
# ------------------------

for file in REMOVE_FILES:
    file_path = ROOT_DIR / file
    if file_path.exists():
        print(f"Removing unnecessary file: {file_path}")
        os.remove(file_path)

print("\n=============================")
print("All operations complete!")
print(f"Files saved under: {ROOT_DIR}")
print("=============================\n")
