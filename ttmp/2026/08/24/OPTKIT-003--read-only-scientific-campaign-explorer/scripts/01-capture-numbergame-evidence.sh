#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
ticket_rel="ttmp/2026/08/24/OPTKIT-003--read-only-scientific-campaign-explorer"
out="$repo_root/$ticket_rel/sources/01-numbergame-evidence.md"
store="$(mktemp -d)"
summary="$(mktemp)"
inspection="$(mktemp)"
verification="$(mktemp)"
trap 'rm -rf "$store" "$summary" "$inspection" "$verification"' EXIT

cd "$repo_root"
GOWORK=off go run ./cmd/optkit demo --store "$store" --reset > "$summary"
campaign_id="$(jq -r .campaign "$summary")"
GOWORK=off go run ./cmd/optkit campaign inspect --store "$store" --id "$campaign_id" --tail 100 > "$inspection"
GOWORK=off go run ./cmd/optkit campaign verify --store "$store" --id "$campaign_id" > "$verification"

{
  printf '# Numbergame Evidence Snapshot\n\n'
  printf 'Generated from Optkit revision `%s`. The temporary local store was deleted after capture.\n\n' "$(git rev-parse HEAD)"
  printf '## Demo summary\n\n```json\n'
  jq 'del(.store_root, .database_path, .artifacts_path)' "$summary"
  printf '```\n\n## Campaign inspection\n\n```json\n'
  cat "$inspection"
  printf '```\n\n## Verification\n\n```json\n'
  cat "$verification"
  printf '```\n'
} > "$out"

printf 'wrote %s\n' "$out"
