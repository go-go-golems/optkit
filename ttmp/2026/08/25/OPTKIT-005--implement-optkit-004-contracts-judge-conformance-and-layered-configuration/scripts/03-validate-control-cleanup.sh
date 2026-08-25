#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"

forbidden='ResourceClaim|ResourceClaims|resource_claims_json|Heartbeat\(|Causation|Correlation|CorrelationID|optionalCorrelation|optionalEvent'
if rg -n "$forbidden" \
  "$OPTKIT/record" "$OPTKIT/scheduler" "$OPTKIT/campaign" \
  "$OPTKIT/query" "$OPTKIT/store" "$OPTKIT/examples" \
  --glob '*.go' --glob '*.sql'; then
  echo "removed control field remains in Optkit source" >&2
  exit 1
fi

cleanup() {
  if [[ -n "${CLEAN_WORKTREE:-}" && -d "${CLEAN_WORKTREE:-}" ]]; then
    git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

CLEAN_WORKTREE=$(mktemp -d)
rmdir "$CLEAN_WORKTREE"
git -C "$OPTKIT" worktree add --detach "$CLEAN_WORKTREE" HEAD >/dev/null
DEMO_STORE=$(mktemp -d)
(
  cd "$CLEAN_WORKTREE"
  GOWORK=off go test ./... -count=1
  GOWORK=off go vet ./...
  GOWORK=off go run ./cmd/optkit demo --store "$DEMO_STORE" --reset >/dev/null
)
rm -rf "$DEMO_STORE"
git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null
CLEAN_WORKTREE=""

(
  cd "$RAG_TTC"
  GOWORK=off go test ./pkg/ttc/optkitcampaign ./cmd/rag-ttc/cmds/experiments/optkitrag -count=1
  GOWORK=off go test -race ./pkg/ttc/optkitcampaign -count=1
)

printf 'OPTKIT_COMMIT=%s\n' "$(git -C "$OPTKIT" rev-parse HEAD)"
printf 'RAG_TTC_COMMIT=%s\n' "$(git -C "$RAG_TTC" rev-parse HEAD)"
printf 'CONTROL_CLEANUP_VALIDATION=PASS\n'
