#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$repo_root/.." && pwd)"
archive="${OPTKIT_SOURCE_ARCHIVE:-$workspace_root/sources/optkit-implementation-source.zip}"
ticket_rel="ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice"
manifest="$repo_root/$ticket_rel/sources/01-optkit-baseline-manifest.txt"

if [[ ! -f "$archive" ]]; then
  printf 'archive not found: %s\n' "$archive" >&2
  exit 1
fi

workdir="$(mktemp -d)"
trap 'rm -rf "$workdir"' EXIT
unzip -q "$archive" -d "$workdir"
source_root="$workdir/optkit"

if [[ "$(sed -n '1p' "$source_root/go.mod")" != "module github.com/go-go-golems/optkit" ]]; then
  printf 'unexpected source module: %s\n' "$(sed -n '1p' "$source_root/go.mod")" >&2
  exit 1
fi

{
  printf 'archive: %s\n' "$archive"
  printf 'archive_sha256: %s\n' "$(sha256sum "$archive" | awk '{print $1}')"
  printf 'source_revision: 1786d1da86c9e03316ed71336bbc993fa30531f0\n'
  printf 'generated_utc: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"
  printf '\nfiles:\n'
  (
    cd "$source_root"
    find . -type f -print0 | sort -z | xargs -0 sha256sum
  )
} > "$manifest"

source_dirs=(
  artifact budget campaign docs episode examples experiment internal local
  measure projection record scheduler space store
)
for dir in "${source_dirs[@]}"; do
  rm -rf "$repo_root/$dir"
  cp -a "$source_root/$dir" "$repo_root/$dir"
done

rm -rf "$repo_root/cmd/XXX" "$repo_root/cmd/optkit"
mkdir -p "$repo_root/cmd"
cp -a "$source_root/cmd/optkit" "$repo_root/cmd/optkit"

cp "$source_root/go.mod" "$repo_root/go.mod"
cp "$source_root/README.md" "$repo_root/README.md"
rm -f "$repo_root/go.sum" "$repo_root/logcopter_generate.go"
rm -rf "$repo_root/pkg"

printf 'Imported Optkit baseline from %s\n' "$archive"
printf 'Manifest: %s\n' "$manifest"
