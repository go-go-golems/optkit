#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
judgekit_root="$workspace_root/judgekit"
rag_ttc_root="$workspace_root/rag-ttc"
expected_work="$workspace_root/go.work"

test "$(cd "$rag_ttc_root" && go env GOWORK)" = "$expected_work"

cd "$judgekit_root"
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run ./...

# The caller may have unrelated unstaged Optkit work. Validate the RAG-TTC
# workspace against a clean detached Optkit worktree without modifying it.
tmp_root="$(mktemp -d)"
clean_optkit="$tmp_root/optkit"
cleanup() {
  git -C "$optkit_root" worktree remove --force "$clean_optkit" >/dev/null 2>&1 || true
  rm -rf "$tmp_root"
}
trap cleanup EXIT
git -C "$optkit_root" worktree add --detach "$clean_optkit" HEAD >/dev/null
cat >"$tmp_root/go.work" <<EOF
go 1.26.6

use (
	$workspace_root/coinvault
	$workspace_root/geppetto
	$workspace_root/judgekit
	$clean_optkit
	$workspace_root/pinocchio
	$workspace_root/rag-ttc
	$workspace_root/ragkit
	$workspace_root/ragopt
	$workspace_root/react-chat
	$workspace_root/sessionstream
)

replace github.com/go-go-golems/judgekit v0.0.0 => $workspace_root/judgekit
replace github.com/go-go-golems/optkit v0.0.0 => $clean_optkit
EOF

cd "$rag_ttc_root"
export GOWORK="$tmp_root/go.work"
go test ./pkg/ttc/judgeinstrument -run 'TestHistorical|TestCacheBypass|TestJudgeFailure' -v -count=1
go test ./... -count=1
go test -race ./... -count=1
go build -buildvcs=false ./...
go vet ./...
golangci-lint run ./...

glazed_lint_bin="$(mktemp)"
go build -buildvcs=false -o "$glazed_lint_bin" github.com/go-go-golems/glazed/cmd/tools/glazed-lint
go vet -vettool="$glazed_lint_bin" ./...
rm -f "$glazed_lint_bin"

git -C "$judgekit_root" diff --check
git -C "$rag_ttc_root" diff --check
git -C "$optkit_root" diff --check

printf 'P7_JUDGEKIT_COMMIT=%s\n' "$(git -C "$judgekit_root" rev-parse HEAD)"
printf 'P7_RAG_TTC_COMMIT=%s\n' "$(git -C "$rag_ttc_root" rev-parse HEAD)"
printf 'P7_VALIDATION=PASS\n'
