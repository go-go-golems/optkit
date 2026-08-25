#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
JUDGEKIT="$WORKSPACE/judgekit"
RAG_TTC="$WORKSPACE/rag-ttc"

cleanup() {
  if [[ -n "${CLEAN_WORKTREE:-}" && -d "${CLEAN_WORKTREE:-}" ]]; then
    git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null 2>&1 || true
  fi
}
trap cleanup EXIT

OPTKIT_VERSION=$(cd "$RAG_TTC" && GOWORK=off go list -m -f '{{.Version}}' github.com/go-go-golems/optkit)
JUDGEKIT_VERSION=$(cd "$RAG_TTC" && GOWORK=off go list -m -f '{{.Version}}' github.com/go-go-golems/judgekit)

[[ "$OPTKIT_VERSION" != "v0.0.0" ]]
[[ "$JUDGEKIT_VERSION" != "v0.0.0" ]]
printf 'OPTKIT_VERSION=%s\n' "$OPTKIT_VERSION"
printf 'JUDGEKIT_VERSION=%s\n' "$JUDGEKIT_VERSION"

CLEAN_WORKTREE=$(mktemp -d)
rmdir "$CLEAN_WORKTREE"
git -C "$OPTKIT" worktree add --detach "$CLEAN_WORKTREE" HEAD >/dev/null
(
  cd "$CLEAN_WORKTREE"
  GOWORK=off go test ./... -count=1
  GOWORK=off go vet ./...
)
git -C "$OPTKIT" worktree remove --force "$CLEAN_WORKTREE" >/dev/null
CLEAN_WORKTREE=""

(
  cd "$JUDGEKIT"
  GOWORK=off go test ./... -count=1
  GOWORK=off go vet ./...
)

(
  cd "$RAG_TTC"
  GOWORK=off go test ./... -count=1
  GOWORK=off go build ./...
  GOWORK=off go vet ./...
  GOWORK=off golangci-lint run
)

printf 'PHASE_R_VALIDATION=PASS\n'
