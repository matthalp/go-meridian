#!/usr/bin/env python3
"""Check that the working tree is ready to be released as a given version.

Usage: check_release.py vX.Y.Z [--notes FILE] [--gorelease PATH]

Run from the repository root. Checks that:
  - the version is well-formed and not already tagged
  - its major version matches the module path in go.mod (/vN suffix)
  - the Version constant in meridian.go matches
  - CHANGELOG.md has a "## [X.Y.Z] - YYYY-MM-DD" section (written to
    --notes, if given, for use as release notes)
  - the bump is right for the API changes since the previous release of the
    same major version, according to gorelease. A change to the exported
    Version constant's value is ignored: gorelease reports any constant value
    change as incompatible, but changing it is the point of a release.

Exits non-zero, with a message per problem, if any check fails.
"""

import argparse
import re
import subprocess
import sys

SEMVER = re.compile(r"^v(0|[1-9]\d*)\.(0|[1-9]\d*)\.(0|[1-9]\d*)$")
IGNORED_INCOMPATIBLE = re.compile(r"^Version: value changed from ")


def run(*args, check=True):
    return subprocess.run(args, check=check, capture_output=True, text=True)


def parse(version):
    m = SEMVER.match(version)
    return tuple(int(x) for x in m.groups()) if m else None


def module_major(go_mod):
    m = re.search(r"^module\s+(\S+)", go_mod, re.M)
    if not m:
        return None, None
    path = m.group(1)
    suffix = re.search(r"/v(\d+)$", path)
    return path, int(suffix.group(1)) if suffix else None


def previous_tag(version):
    """Return the highest existing tag with the same major version, below version."""
    tags = run("git", "tag", "--list", f"v{version[0]}.*").stdout.split()
    older = sorted((parse(t), t) for t in tags if parse(t) and parse(t) < version)
    return older[-1][1] if older else None


def changelog_section(text, bare_version):
    heading = re.compile(rf"^## \[{re.escape(bare_version)}\] - \d{{4}}-\d{{2}}-\d{{2}}\s*$", re.M)
    m = heading.search(text)
    if not m:
        return None
    nxt = re.compile(r"^## \[", re.M).search(text, m.end())
    return text[m.end():nxt.start() if nxt else len(text)].strip()


def gorelease_sections(report):
    """Collect the lines under each '## <section>' heading of a gorelease report."""
    sections, current = {}, None
    for line in report.splitlines():
        if line.startswith("## "):
            current = line[3:].strip()
            sections.setdefault(current, [])
        elif line.startswith("# "):
            current = None
        elif current and line.strip():
            sections[current].append(line.strip())
    return sections


def check_api(version_str, version, base, gorelease):
    """Return a list of problems with the version bump relative to base."""
    proc = run(gorelease, f"-base={base}", f"-version={version_str}", check=False)
    report = proc.stdout + proc.stderr
    print(f"gorelease report (base {base}):\n{report}")

    sections = gorelease_sections(report)
    # Changes are reported once per package, so a change to Time[TZ] repeats
    # for every timezone package's alias; keep each distinct line once.
    raw_incompatible = list(dict.fromkeys(sections.pop("incompatible changes", [])))
    incompatible = [l for l in raw_incompatible if not IGNORED_INCOMPATIBLE.match(l)]
    compatible = list(dict.fromkeys(sections.pop("compatible changes", [])))
    problems = [f"gorelease: {name}: {l}" for name, lines in sections.items() for l in lines]

    base_v = parse(base)
    if incompatible:
        problems.append(
            "incompatible API changes need a new major version (and module path); "
            "found: " + "; ".join(incompatible)
        )
    elif compatible and version[:2] == base_v[:2]:
        problems.append(
            f"{version_str} is a patch release but there are new APIs; use a minor version: " + "; ".join(compatible)
        )
    if proc.returncode != 0 and not raw_incompatible and not problems:
        # gorelease failed for a reason its change sections don't explain,
        # e.g. uncommitted changes in the working tree.
        last = report.strip().splitlines()[-1] if report.strip() else f"exit status {proc.returncode}"
        problems.append(f"gorelease failed: {last}")
    return problems


def main():
    ap = argparse.ArgumentParser()
    ap.add_argument("version")
    ap.add_argument("--notes", help="write the CHANGELOG section for this version to this file")
    ap.add_argument("--gorelease", default="gorelease", help="path to the gorelease binary")
    args = ap.parse_args()

    v_str = args.version
    version = parse(v_str)
    if not version:
        sys.exit(f"error: {v_str!r} is not of the form vMAJOR.MINOR.PATCH")
    bare = v_str[1:]
    problems = []

    if run("git", "tag", "--list", v_str).stdout.strip():
        problems.append(f"tag {v_str} already exists")

    path, suffix = module_major(open("go.mod").read())
    want_suffix = version[0] if version[0] >= 2 else None
    if suffix != want_suffix:
        problems.append(f"module path {path} does not match major version {version[0]}")

    const = re.search(r'^const Version = "([^"]*)"', open("meridian.go").read(), re.M)
    if const and const.group(1) != bare:
        problems.append(f'meridian.go has Version = "{const.group(1)}", want "{bare}"')

    notes = changelog_section(open("CHANGELOG.md").read(), bare)
    if notes is None:
        problems.append(f"CHANGELOG.md has no '## [{bare}] - YYYY-MM-DD' section")
    elif args.notes:
        with open(args.notes, "w") as f:
            f.write(notes + "\n")

    base = previous_tag(version)
    if base is None:
        print(f"No earlier v{version[0]}.x tag; skipping the API compatibility check.")
    else:
        problems += check_api(v_str, version, base, args.gorelease)

    if problems:
        for p in problems:
            print(f"::error::{p}")
        sys.exit(1)
    print(f"{v_str} is ready to release.")


if __name__ == "__main__":
    main()
