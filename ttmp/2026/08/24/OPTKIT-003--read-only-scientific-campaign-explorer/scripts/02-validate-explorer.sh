#!/usr/bin/env bash
set -euo pipefail

repo_root="$(git rev-parse --show-toplevel)"
cd "$repo_root"

run() {
  printf '\n$ %s\n' "$*"
  "$@"
}

printf '# OPTKIT-003 explorer validation\n'
printf 'revision: %s\n' "$(git rev-parse HEAD)"
printf 'validated_utc: %s\n' "$(date -u +%Y-%m-%dT%H:%M:%SZ)"

run make ci-check
run make race
run make lint

store="$(mktemp -d)"
summary="$(mktemp)"
server_log="$(mktemp)"
port="${OPTKIT_EXPLORER_TEST_PORT:-18081}"
server_pid=""
cleanup() {
  if [[ -n "$server_pid" ]]; then kill "$server_pid" 2>/dev/null || true; fi
  rm -rf "$store" "$summary" "$server_log"
}
trap cleanup EXIT

GOWORK=off go run ./cmd/optkit demo --store "$store" --reset > "$summary"
lsof-who -p "$port" -k >/dev/null 2>&1 || true
campaign_id="$(jq -r .campaign "$summary")"
GOWORK=off go run ./cmd/optkit serve --store "$store" --listen "127.0.0.1:$port" > "$server_log" 2>&1 &
server_pid="$!"

for _ in $(seq 1 60); do
  if curl -fsS "http://127.0.0.1:$port/api/v1/health" >/dev/null; then break; fi
  sleep 0.25
done

run curl -fsS "http://127.0.0.1:$port/api/v1/health"
run curl -fsS "http://127.0.0.1:$port/api/v1/campaigns"
run curl -fsS "http://127.0.0.1:$port/api/v1/campaigns/$campaign_id"
run curl -fsS "http://127.0.0.1:$port/api/v1/campaigns/$campaign_id/events?limit=2"
run curl -fsS "http://127.0.0.1:$port/"

post_status="$(curl -sS -o /dev/null -w '%{http_code}' -X POST "http://127.0.0.1:$port/api/v1/campaigns")"
if [[ "$post_status" != "405" ]]; then
  printf 'POST /api/v1/campaigns status=%s, want 405\n' "$post_status" >&2
  exit 1
fi
printf '\nPOST_MUTATION_STATUS=%s\n' "$post_status"
printf 'CAMPAIGN=%s\n' "$campaign_id"
printf 'EXPLORER_VALIDATION=PASS\n'
