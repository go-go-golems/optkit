#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
rag_ttc_root="$workspace_root/rag-ttc"

cd "$rag_ttc_root"

GOWORK=off go test ./pkg/ttc/retrievaleval ./cmd/rag-ttc/cmds/experiments/answerquality -count=1
GOWORK=off go test -race ./pkg/ttc/retrievaleval ./cmd/rag-ttc/cmds/experiments/answerquality -count=1
GOWORK=off go test ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run
make glazed-lint
git diff --check

printf 'P4_COMMIT=%s\n' "$(git rev-parse HEAD)"
printf 'P4_VALIDATION=PASS\n'
