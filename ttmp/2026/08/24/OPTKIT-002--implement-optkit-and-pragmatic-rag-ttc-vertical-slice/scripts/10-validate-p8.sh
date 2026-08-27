#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
ragkit_root="$workspace_root/ragkit"
judgekit_root="$workspace_root/judgekit"
ragopt_root="$workspace_root/ragopt"
rag_ttc_root="$workspace_root/rag-ttc"
coinvault_root="$workspace_root/coinvault"

if rg -n 'ToolRegistryFactory|buildProviderToolRegistry' "$rag_ttc_root" --glob '*.go'; then
  echo 'superseded outer customer orchestration still has an active caller' >&2
  exit 1
fi

# Candidate/eval orchestration is retained only in the frozen I5 command until
# full answer/gate/report parity exists. No other active RAG-TTC source may add
# a new RagOpt candidate/eval/gate/policy/compare/report dependency.
unexpected_ragopt="$(rg -l '"github.com/go-go-golems/ragopt/pkg/(candidate|eval|gate|policy|compare|report)' "$rag_ttc_root" --glob '*.go' | grep -v '/cmd/rag-ttc/cmds/tooleval/ragopt' || true)"
test -z "$unexpected_ragopt"

# Use committed Optkit state so unrelated local edits cannot affect boundary or
# cross-repository validation.
tmp_root="$(mktemp -d)"
clean_optkit="$tmp_root/optkit"
cleanup() {
  git -C "$optkit_root" worktree remove --force "$clean_optkit" >/dev/null 2>&1 || true
  rm -rf "$tmp_root"
}
trap cleanup EXIT
git -C "$optkit_root" worktree add --detach "$clean_optkit" HEAD >/dev/null

cd "$clean_optkit"
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./... -count=1
GOWORK=off go vet ./...
GOWORK=off golangci-lint run ./...

cd "$ragkit_root"
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run ./...

cd "$judgekit_root"
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run ./...

cd "$ragopt_root"
GOWORK=off go test ./... -count=1
GOWORK=off go test -race ./... -count=1
GOWORK=off go build ./...
GOWORK=off go vet ./...
GOWORK=off golangci-lint run ./...

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
go test ./pkg/ttc/search -run TestCrossProductSemanticFixtureCapturesCurrentSearchAndEvidenceLaws -count=1
go test ./pkg/ttc/optkitcampaign ./pkg/ttc/judgeinstrument ./internal/customer/realruntime ./internal/customer/webchatcmd -count=1
go test ./... -count=1
go test -race ./... -count=1
go build -buildvcs=false ./...
go vet ./...
golangci-lint run ./...

glazed_lint_bin="$(mktemp)"
go build -buildvcs=false -o "$glazed_lint_bin" github.com/go-go-golems/glazed/cmd/tools/glazed-lint
go vet -vettool="$glazed_lint_bin" ./...
rm -f "$glazed_lint_bin"

cd "$coinvault_root"
GOWORK=off go test ./cmd/coinvault/cmds -run Ragopt -count=1

for repository in "$optkit_root" "$ragkit_root" "$judgekit_root" "$ragopt_root" "$rag_ttc_root"; do
  git -C "$repository" diff --check
done

printf 'P8_OPTKIT_COMMIT=%s\n' "$(git -C "$optkit_root" rev-parse HEAD)"
printf 'P8_RAGKIT_COMMIT=%s\n' "$(git -C "$ragkit_root" rev-parse HEAD)"
printf 'P8_JUDGEKIT_COMMIT=%s\n' "$(git -C "$judgekit_root" rev-parse HEAD)"
printf 'P8_RAG_TTC_COMMIT=%s\n' "$(git -C "$rag_ttc_root" rev-parse HEAD)"
printf 'P8_RAGOPT_RETAINED_COMMIT=%s\n' "$(git -C "$ragopt_root" rev-parse HEAD)"
printf 'P8_VALIDATION=PASS\n'
