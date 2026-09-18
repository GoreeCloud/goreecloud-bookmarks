#!/usr/bin/env python3
import json
import zipfile
from pathlib import Path

ROOT = Path(__file__).resolve().parents[3]
CLIENT_DIR = ROOT / "clients" / "firefox"
EXCLUDE_NAMES = {"LICENSE", ".source-baseline", "release-state.json"}
EXCLUDE_SUFFIXES = {".md", ".py", ".pyc"}
EXCLUDE_PARTS = {"scripts", "tests", "__pycache__", ".git"}

def should_include(path: Path) -> bool:
    rel = path.relative_to(CLIENT_DIR)
    if any(part in EXCLUDE_PARTS for part in rel.parts):
        return False
    if path.name in EXCLUDE_NAMES:
        return False
    if path.suffix.lower() in EXCLUDE_SUFFIXES:
        return False
    return path.is_file()

def main() -> None:
    manifest_path = CLIENT_DIR / "manifest.json"
    manifest = json.loads(manifest_path.read_text(encoding="utf-8"))
    version = manifest["version"]
    output_dir = ROOT / "dist"
    output_dir.mkdir(parents=True, exist_ok=True)
    output = output_dir / f"goreecloud-bookmarks-{version}.xpi"
    files = sorted(path for path in CLIENT_DIR.rglob("*") if should_include(path))
    with zipfile.ZipFile(output, "w", compression=zipfile.ZIP_DEFLATED, compresslevel=9) as archive:
        for path in files:
            rel = path.relative_to(CLIENT_DIR).as_posix()
            info = zipfile.ZipInfo(rel, date_time=(1980, 1, 1, 0, 0, 0))
            info.compress_type = zipfile.ZIP_DEFLATED
            info.external_attr = 0o644 << 16
            archive.writestr(info, path.read_bytes())
    with zipfile.ZipFile(output) as archive:
        bad = archive.testzip()
        if bad:
            raise SystemExit(f"Package integrity failure: {bad}")
        if "manifest.json" not in archive.namelist():
            raise SystemExit("Package does not contain manifest.json at archive root")
    print(output.relative_to(ROOT))

if __name__ == "__main__":
    main()
