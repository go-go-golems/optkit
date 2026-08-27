#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"
BASELINE="$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v1.yaml"
CHALLENGER="$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-challenger-v1.yaml"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
BIN="$TMP/rag-ttc"
STORE="$TMP/store"

cd "$RAG_TTC"
GOWORK=off go build -o "$BIN" ./cmd/rag-ttc
GOWORK=off go test ./pkg/ttc/experimentworkbench ./pkg/ttc/optkitcampaign ./pkg/ttc/optimization ./cmd/rag-ttc/cmds/experiments/optkitrag -count=1
GOWORK=off go test -race ./pkg/ttc/experimentworkbench ./pkg/ttc/optkitcampaign ./cmd/rag-ttc/cmds/experiments/optkitrag -count=1

"$BIN" experiment optkit-rag config validate --manifest "$BASELINE" --format json > "$TMP/validate.json"
jq -e 'length == 1 and .[0].status == "valid" and .[0].arms == 2 and .[0].cases == 3 and .[0].episodes == 6 and (.[0].graph_ids | length == 2)' "$TMP/validate.json" >/dev/null
printf 'MANIFEST_VALIDATION=PASS\n'

"$BIN" experiment optkit-rag config inspect --manifest "$BASELINE" --format json > "$TMP/inspect.json"
jq -e 'length == 24 and ([.[].arm] | unique | length == 2) and ([.[].layer] | unique | length == 12)' "$TMP/inspect.json" >/dev/null
printf 'PER_ARM_GRAPH_INSPECTION=PASS\n'

"$BIN" experiment optkit-rag config diff \
  --before "$BASELINE" --before-arm limit-2 \
  --after "$CHALLENGER" --after-arm limit-3 \
  --format json > "$TMP/diff.json"
jq -e '([.[] | select(.changed == true)] | length == 1) and ([.[] | select(.changed == true)][0].layer == "retrieval")' "$TMP/diff.json" >/dev/null
printf 'DIRECT_CONFIG_DIFF=PASS\n'

"$BIN" experiment optkit-rag config plan \
  --before "$BASELINE" --before-arm limit-2 \
  --after "$CHALLENGER" --after-arm limit-3 \
  --format json > "$TMP/plan.json"
jq -e '([.[] | select(.action == "reuse")] | length == 5) and ([.[] | select(.action == "recompute")] | length == 7) and ([.[] | select(.reason == "direct_change")][0].layer == "retrieval")' "$TMP/plan.json" >/dev/null
printf 'INVALIDATION_PLAN=PASS\n'

"$BIN" experiment optkit-rag campaign dry-run --manifest "$BASELINE" --store "$STORE" --format json > "$TMP/dry-run.json"
jq -e 'length == 1 and .[0].mutation == false and .[0].episodes == 6 and .[0].retrieval_query_budget == 6 and .[0].retrieval_result_budget == 600' "$TMP/dry-run.json" >/dev/null
test ! -e "$STORE"
printf 'DRY_RUN_NO_MUTATION=PASS\n'

"$BIN" experiment optkit-rag campaign run --manifest "$BASELINE" --store "$STORE" --format json > "$TMP/run.json"
CAMPAIGN=$(jq -er '.[0].campaign' "$TMP/run.json")
TRIAL=$(jq -er '.[0].trial' "$TMP/run.json")
VERSION=$(jq -er '.[0].version' "$TMP/run.json")
jq -e 'length == 1 and .[0].status == "completed" and .[0].completed == 6 and .[0].failed_terminal == 0 and .[0].budget_violated == false' "$TMP/run.json" >/dev/null
printf 'MANIFEST_CAMPAIGN_RUN=PASS\n'

"$BIN" experiment optkit-rag campaign status --store "$STORE" --campaign "$CAMPAIGN" --format json > "$TMP/status.json"
jq -e --arg campaign "$CAMPAIGN" --arg trial "$TRIAL" '.[0].campaign == $campaign and .[0].trial == $trial and .[0].status == "completed"' "$TMP/status.json" >/dev/null
printf 'CAMPAIGN_STATUS=PASS\n'

"$BIN" experiment optkit-rag campaign verify --store "$STORE" --campaign "$CAMPAIGN" --format json > "$TMP/verify.json"
jq -e '.[0].journal_verified == true and .[0].direct_payloads_verified == true and .[0].events > 0 and .[0].unique_direct_payloads > 0' "$TMP/verify.json" >/dev/null
printf 'CAMPAIGN_CUSTODY=PASS\n'

"$BIN" experiment optkit-rag campaign resume --store "$STORE" --campaign "$CAMPAIGN" --format json > "$TMP/resume.json"
jq -e --arg campaign "$CAMPAIGN" --arg trial "$TRIAL" --argjson version "$VERSION" '.[0].campaign == $campaign and .[0].trial == $trial and .[0].version == $version and .[0].completed == 6' "$TMP/resume.json" >/dev/null
printf 'IDEMPOTENT_RESUME=PASS\n'

"$BIN" experiment optkit-rag campaign run --help > "$TMP/help.txt"
for flag in manifest store reset; do grep -q -- "--$flag" "$TMP/help.txt"; done
"$BIN" experiment optkit-rag campaign run --help --long-help > "$TMP/long-help.txt"
for flag in format output-fields max-output-rows; do grep -q -- "--$flag" "$TMP/long-help.txt"; done
printf 'GLAZED_HELP_SURFACE=PASS\n'

printf 'CAMPAIGN=%s\n' "$CAMPAIGN"
printf 'TRIAL=%s\n' "$TRIAL"
printf 'OPTKIT_006_WORKBENCH_VALIDATION=PASS\n'
