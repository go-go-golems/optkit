#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
rag_ttc_root="$workspace_root/rag-ttc"

cd "$rag_ttc_root"

focused=(
  ./pkg/ttc/customerapp
  ./pkg/ttc/search
  ./internal/customer/ragsearch
  ./internal/customer/realruntime
  ./internal/customer/webchatcmd
)
GOWORK=off go test "${focused[@]}" -count=1
GOWORK=off go test -race "${focused[@]}" -count=1
GOWORK=off go test ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run
make glazed-lint
git diff --check

printf 'P5_COMMIT=%s\n' "$(git rev-parse HEAD)"
printf 'P5_VALIDATION=PASS\n'
