#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"

"$SCRIPT_DIR/02-validate-model-doc.py"
"$SCRIPT_DIR/03-validate-control-cleanup.sh"
"$SCRIPT_DIR/04-validate-phase0.sh"
"$SCRIPT_DIR/05-validate-phase1.sh"
"$SCRIPT_DIR/06-validate-phase2.sh"

cleanup() {
  if [[ -n "${CLEAN_WORKTREE:-}" && -d "${CLEAN_WORKTREE:-}" ]]; then
    git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

CLEAN_WORKTREE=$(mktemp -d)
rmdir "$CLEAN_WORKTREE"
git -C "$OPTKIT" worktree add --detach "$CLEAN_WORKTREE" HEAD >/dev/null
(
  cd "$CLEAN_WORKTREE"
  GOWORK=off go test -race ./... -count=1
  GOWORK=off go build ./...
  GOWORK=off go vet ./...
)
git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null
CLEAN_WORKTREE=""

(
  cd "$RAG_TTC"
  GOWORK=off go test ./... -count=1
  GOWORK=off go test -race ./pkg/ttc/... ./internal/admin/chatserver/... -count=1
  GOWORK=off go build ./...
  GOWORK=off go vet ./...
  GOWORK=off golangci-lint run
)

printf 'FINAL_OPTKIT_COMMIT=%s\n' "$(git -C "$OPTKIT" rev-parse HEAD)"
printf 'FINAL_RAG_TTC_COMMIT=%s\n' "$(git -C "$RAG_TTC" rev-parse HEAD)"
printf 'OPTKIT_005_FINAL_VALIDATION=PASS\n'
