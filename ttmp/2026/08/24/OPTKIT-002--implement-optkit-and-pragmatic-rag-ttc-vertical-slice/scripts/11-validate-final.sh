#!/usr/bin/env bash
set -euo pipefail
set -x

optkit_root="$(git rev-parse --show-toplevel)"
workspace_root="$(cd "$optkit_root/.." && pwd)"
ticket_root="$optkit_root/ttmp/2026/08/24/OPTKIT-002--implement-optkit-and-pragmatic-rag-ttc-vertical-slice"
rag_ttc_root="$workspace_root/rag-ttc"

for phase in 0 3 4 5 6 7 8; do
  case "$phase" in
    0) evidence="$ticket_root/sources/02-p0-validation.txt" ;;
    3) evidence="$ticket_root/sources/04-p3-validation.txt" ;;
    4) evidence="$ticket_root/sources/05-p4-validation.txt" ;;
    5) evidence="$ticket_root/sources/06-p5-validation.txt" ;;
    6) evidence="$ticket_root/sources/07-p6-validation.txt" ;;
    7) evidence="$ticket_root/sources/08-p7-validation.txt" ;;
    8) evidence="$ticket_root/sources/10-p8-validation.txt" ;;
  esac
  grep -q "P${phase}_VALIDATION=PASS" "$evidence"
done

"$ticket_root/scripts/03-sync-rag-semantic-fixture.sh" --check
"$ticket_root/scripts/10-validate-p8.sh"

tmp_root="$(mktemp -d)"
clean_optkit="$tmp_root/optkit"
store_root="$tmp_root/store"
output_root="$tmp_root/output"
cleanup() {
  git -C "$optkit_root" worktree remove --force "$clean_optkit" >/dev/null 2>&1 || true
  rm -rf "$tmp_root"
}
trap cleanup EXIT
mkdir -p "$output_root"
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
go run ./cmd/rag-ttc experiment optkit-rag run --store "$store_root" --reset --format json >"$output_root/run.json"
campaign_id="$(python3 -c 'import json,sys; rows=json.load(open(sys.argv[1])); assert len(rows)==1; value=rows[0]; assert value["status"]=="completed"; assert value["completed"]==6; assert abs(value["paired_delta"]-(1/6))<1e-12; print(value["campaign"])' "$output_root/run.json")"
go run ./cmd/rag-ttc experiment optkit-rag run --store "$store_root" --campaign "$campaign_id" --format json >"$output_root/resume.json"
go run ./cmd/rag-ttc experiment optkit-rag inspect --store "$store_root" --campaign "$campaign_id" --format json >"$output_root/inspect.json"
cmp "$output_root/run.json" "$output_root/resume.json"
cmp "$output_root/resume.json" "$output_root/inspect.json"

cd "$clean_optkit"
GOWORK=off CGO_ENABLED=0 go test ./... -count=1
GOWORK=off CGO_ENABLED=0 go build ./...

printf 'FINAL_CAMPAIGN=%s\n' "$campaign_id"
printf 'FINAL_OPTKIT_COMMIT=%s\n' "$(git -C "$optkit_root" rev-parse HEAD)"
printf 'FINAL_RAG_TTC_COMMIT=%s\n' "$(git -C "$rag_ttc_root" rev-parse HEAD)"
printf 'FINAL_VALIDATION=PASS\n'
