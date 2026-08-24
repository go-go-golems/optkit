#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

printf '# OPTKIT-002 P0 validation\n\n'
printf 'revision: %s\n' "$(git rev-parse HEAD)"
printf 'validated_utc: %s\n\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"

run() {
  printf '\n$ %s\n' "$*"
  "$@"
}

run make ci-check
run make race
run make lint

store="$(mktemp -d)"
trap 'rm -rf "$store"' EXIT
summary="$store-summary.json"
printf '\n$ GOWORK=off go run ./cmd/optkit demo --store %s --reset\n' "$store"
GOWORK=off go run ./cmd/optkit demo --store "$store" --reset | tee "$summary"
campaign_id="$(python3 -c 'import json,sys; print(json.load(open(sys.argv[1]))["campaign"])' "$summary")"
run env GOWORK=off go run ./cmd/optkit campaign inspect --store "$store" --id "$campaign_id"
run env GOWORK=off go run ./cmd/optkit campaign verify --store "$store" --id "$campaign_id"

printf '\nP0_VALIDATION=PASS\n'
