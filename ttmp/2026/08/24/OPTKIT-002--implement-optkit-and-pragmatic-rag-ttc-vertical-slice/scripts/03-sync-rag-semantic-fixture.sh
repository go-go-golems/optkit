#!/usr/bin/env bash
set -euo pipefail

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
ticket_rel="ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice"
source_fixture="$optkit_root/$ticket_rel/sources/rag-semantic-fixture-v1.json"
targets=(
  "$workspace_root/rag-ttc/pkg/ttc/search/testdata/rag-semantic-fixture-v1.json"
  "$workspace_root/coinvault/internal/knowledge/testdata/rag-semantic-fixture-v1.json"
)

if [[ ! -f "$source_fixture" ]]; then
  printf 'canonical fixture not found: %s\n' "$source_fixture" >&2
  exit 1
fi

case "${1:-sync}" in
  sync)
    for target in "${targets[@]}"; do
      mkdir -p "$(dirname "$target")"
      cp "$source_fixture" "$target"
      printf 'synced %s\n' "$target"
    done
    ;;
  --check)
    for target in "${targets[@]}"; do
      if ! cmp -s "$source_fixture" "$target"; then
        printf 'fixture drift: %s\n' "$target" >&2
        exit 1
      fi
      printf 'verified %s\n' "$target"
    done
    ;;
  *)
    printf 'usage: %s [sync|--check]\n' "$0" >&2
    exit 2
    ;;
esac

printf 'sha256=%s\n' "$(sha256sum "$source_fixture" | awk '{print $1}')"
