#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
TICKET=$(dirname "$SCRIPT_DIR")
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
WORKSPACE=$(dirname "$OPTKIT")
RAG_TTC="$WORKSPACE/rag-ttc"
FIXTURES="$TICKET/sources/api-fixtures"
TMP=$(mktemp -d)
SERVER_PID=""
cleanup() {
  if [[ -n "$SERVER_PID" ]]; then kill "$SERVER_PID" >/dev/null 2>&1 || true; wait "$SERVER_PID" >/dev/null 2>&1 || true; fi
  rm -rf "$TMP"
}
trap cleanup EXIT
mkdir -p "$FIXTURES"

BIN="$TMP/rag-ttc"
STORE="$TMP/store"
MANIFEST="$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v1.yaml"
PORT=$(python3 - <<'PY'
import socket
s=socket.socket(); s.bind(('127.0.0.1',0)); print(s.getsockname()[1]); s.close()
PY
)
BASE="http://127.0.0.1:$PORT"

cd "$RAG_TTC"
GOWORK=off go build -o "$BIN" ./cmd/rag-ttc
GOWORK=off go test ./pkg/ttc/specialistapi ./pkg/ttc/experimentworkbench ./cmd/rag-ttc/cmds/experiments/optkitrag -count=1
GOWORK=off go test -race ./pkg/ttc/specialistapi -count=1

"$BIN" experiment optkit-rag campaign run --manifest "$MANIFEST" --store "$STORE" --format json > "$TMP/run.json"
CAMPAIGN=$(jq -er '.[0].campaign' "$TMP/run.json")
"$BIN" experiment optkit-rag campaign serve --store "$STORE" --listen "127.0.0.1:$PORT" >"$TMP/server.out" 2>"$TMP/server.err" &
SERVER_PID=$!
for _ in $(seq 1 100); do
  if curl -fsS "$BASE/api/rag/v1/health" > "$TMP/health.json" 2>/dev/null; then break; fi
  sleep 0.05
done
jq -e '.status == "ok" and .read_only == true' "$TMP/health.json" >/dev/null
printf 'SPECIALIST_HEALTH=PASS\n'

curl -fsS -D "$TMP/cockpit.headers" "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cockpit" > "$FIXTURES/01-cockpit.json"
jq -e '.schema == "rag-ttc.campaign-cockpit/v1" and .integrity.journal_verified == true and (.arms|length)==2 and (.cases|length)==3 and .episodes.completed==6' "$FIXTURES/01-cockpit.json" >/dev/null
grep -qi '^etag:' "$TMP/cockpit.headers"
grep -qi '^cache-control: no-store' "$TMP/cockpit.headers"
printf 'COCKPIT_API=PASS\n'

curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/comparisons/limit-1/limit-2" > "$FIXTURES/02-comparison.json"
jq -e '.schema == "rag-ttc.campaign-comparison/v1" and (.cases|length)==3 and (.config_diff.layers|length)==12 and (.invalidation_plan.steps|length)==12 and (.metrics|length)==1' "$FIXTURES/02-comparison.json" >/dev/null
EPISODE=$(jq -er '.cases[0].treatment_episode' "$FIXTURES/02-comparison.json")
printf 'COMPARISON_API=PASS\n'

curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cases?baseline=limit-1&treatment=limit-2&limit=2" > "$FIXTURES/03-cases-page-1.json"
CURSOR=$(jq -er '.next_cursor' "$FIXTURES/03-cases-page-1.json")
curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cases?baseline=limit-1&treatment=limit-2&limit=2&after=$CURSOR" > "$FIXTURES/04-cases-page-2.json"
jq -e '(.cases|length)==2 and .has_more==true' "$FIXTURES/03-cases-page-1.json" >/dev/null
jq -e '(.cases|length)==1 and .has_more==false' "$FIXTURES/04-cases-page-2.json" >/dev/null
printf 'CASE_PAGINATION=PASS\n'

curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/episodes/$EPISODE/pipeline" > "$FIXTURES/05-pipeline.json"
jq -e '.schema == "rag-ttc.episode-pipeline/v1" and (.stages|length)>0 and .episode != ""' "$FIXTURES/05-pipeline.json" >/dev/null
printf 'PIPELINE_API=PASS\n'

curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/provenance/episode/$EPISODE" > "$FIXTURES/06-provenance.json"
jq -e '.schema == "rag-ttc.episode-provenance/v1" and .graph.id != "" and .snapshot.id != "" and (.artifact_edges|length)>0' "$FIXTURES/06-provenance.json" >/dev/null
printf 'PROVENANCE_API=PASS\n'

STATUS=$(curl -sS -o "$TMP/mutation.txt" -w '%{http_code}' -X POST "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cockpit")
test "$STATUS" = 405
STATUS=$(curl -sS -o "$TMP/bad-limit.json" -w '%{http_code}' "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cases?baseline=limit-1&treatment=limit-2&limit=101")
test "$STATUS" = 400
jq -e '.error.code == "invalid_page" or .error.code == "invalid_limit"' "$TMP/bad-limit.json" >/dev/null
printf 'READ_ONLY_AND_ERRORS=PASS\n'

printf 'CAMPAIGN=%s\n' "$CAMPAIGN"
printf 'EPISODE=%s\n' "$EPISODE"
printf 'OPTKIT_007_SPECIALIST_API_VALIDATION=PASS\n'
