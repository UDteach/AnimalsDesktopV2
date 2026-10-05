#!/usr/bin/env bash
set -euo pipefail

usage() {
  cat <<'EOF'
Usage: scripts/clean_appledouble.sh [--check] [--include-git]

Deletes macOS AppleDouble sidecars (._*) and .DS_Store files from the repo.
Delete mode includes the Git common directory so transient sidecars cannot
break git commands, including when running from a linked worktree.
Check mode scans the working tree by default; use --include-git for a full scan.
EOF
}

mode="delete"
include_git="auto"
while [[ $# -gt 0 ]]; do
  case "$1" in
    --check)
      mode="check"
      ;;
    --include-git)
      include_git="yes"
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      usage >&2
      exit 2
      ;;
  esac
  shift
done

if [[ "$include_git" == "auto" ]]; then
  if [[ "$mode" == "delete" ]]; then
    include_git="yes"
  else
    include_git="no"
  fi
fi

root="$(git rev-parse --show-toplevel)"
git_common_dir="$(git rev-parse --git-common-dir)"
case "$git_common_dir" in
  /*) ;;
  *) git_common_dir="$root/$git_common_dir" ;;
esac
tmp="$(mktemp)"
trap 'rm -f "$tmp"' EXIT
if [[ "$include_git" == "yes" ]]; then
  scan_paths=("$root")
  if [[ -d "$git_common_dir" ]]; then
    scan_paths+=("$git_common_dir")
  fi
  find "${scan_paths[@]}" \
    \( -name '._*' -o -name '.DS_Store' \) \
    -type f \
    -print | sort -u >"$tmp"
else
  find "$root" \
    \( -path "$root/.git" -o -path "$root/.git/*" \) -prune \
    -o \( -name '._*' -o -name '.DS_Store' \) \
    -type f \
    -print >"$tmp"
fi

if [[ ! -s "$tmp" ]]; then
  if [[ "$include_git" == "yes" ]]; then
    echo "No AppleDouble or .DS_Store files found."
  else
    echo "No AppleDouble or .DS_Store files found in the working tree."
  fi
  exit 0
fi

cat "$tmp"
if [[ "$mode" == "check" ]]; then
  echo "AppleDouble or .DS_Store files found." >&2
  exit 1
fi

count=0
while IFS= read -r file; do
  rm -f -- "$file"
  count=$((count + 1))
done <"$tmp"
echo "Deleted $count AppleDouble/.DS_Store file(s)."
