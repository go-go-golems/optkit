#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
rag_ttc_root="$workspace_root/rag-ttc"
expected_work="$workspace_root/go.work"

cd "$rag_ttc_root"
test "$(go env GOWORK)" = "$expected_work"

go test ./pkg/ttc/search ./pkg/ttc/optkitcampaign ./cmd/rag-ttc/cmds/experiments/... -count=1
go test -race ./pkg/ttc/search ./pkg/ttc/optkitcampaign ./cmd/rag-ttc/cmds/experiments/... -count=1
go test ./internal/admin/chatserver -run TestWebSocketHeartbeatTimeoutAndServerCloseAreDeterministic -count=20
go test ./... -count=1
go test -race ./... -count=1
go build -buildvcs=false ./...
go vet ./...
golangci-lint run ./...

glazed_lint_bin="$(mktemp)"
go build -buildvcs=false -o "$glazed_lint_bin" github.com/go-go-golems/glazed/cmd/tools/glazed-lint
go vet -vettool="$glazed_lint_bin" ./...
rm -f "$glazed_lint_bin"

store_root="$(mktemp -d)"
output_root="$(mktemp -d)"
trap 'rm -rf "$store_root" "$output_root"' EXIT

go run ./cmd/rag-ttc experiment optkit-rag run --store "$store_root" --reset --format json > "$output_root/run.json"
campaign_id="$(python3 -c 'import json,sys; rows=json.load(open(sys.argv[1])); assert len(rows) == 1; value=rows[0]; assert value["status"] == "completed"; assert value["completed"] == 6; assert abs(value["paired_delta"] - (1/6)) < 1e-12; print(value["campaign"])' "$output_root/run.json")"
go run ./cmd/rag-ttc experiment optkit-rag run --store "$store_root" --campaign "$campaign_id" --format json > "$output_root/resume.json"
go run ./cmd/rag-ttc experiment optkit-rag inspect --store "$store_root" --campaign "$campaign_id" --format json > "$output_root/inspect.json"
cmp "$output_root/run.json" "$output_root/resume.json"
cmp "$output_root/resume.json" "$output_root/inspect.json"

cd "$optkit_root"
go test ./... -count=1
go test -race ./... -count=1
go vet ./...
golangci-lint run ./...
git diff --check

printf 'P6_RAG_TTC_COMMIT=%s\n' "$(git -C "$rag_ttc_root" rev-parse HEAD)"
printf 'P6_OPTKIT_REGISTRY_COMMIT=%s\n' "$(git rev-parse b8e233e86c5f01144d4d12c408200bd74ff48253)"
printf 'P6_CAMPAIGN=%s\n' "$campaign_id"
printf 'P6_VALIDATION=PASS\n'
