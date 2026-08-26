#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
RAG_TTC="$WORKSPACE/rag-ttc"
TMP=$(mktemp -d)
SERVER_PID=""
cleanup() {
  if [[ -n "$SERVER_PID" ]]; then
    kill "$SERVER_PID" 2>/dev/null || true
    wait "$SERVER_PID" 2>/dev/null || true
  fi
  rm -rf "$TMP"
}
trap cleanup EXIT
STORE="$TMP/store"
TOKEN="live-smoke-token"
ACTOR="actor:live-smoke"
PORT=$(python3 - <<'PY'
import socket
s=socket.socket(); s.bind(('127.0.0.1',0)); print(s.getsockname()[1]); s.close()
PY
)
BASE="http://127.0.0.1:$PORT"

(
  cd "$RAG_TTC"
  GOWORK=off go build -o "$TMP/rag-ttc" ./cmd/rag-ttc
  GOWORK=off go run "$SCRIPT_DIR/setup-running-campaign/main.go" --store "$STORE" >"$TMP/campaign.txt"
)
CAMPAIGN=$(tr -d '\n' <"$TMP/campaign.txt")
"$TMP/rag-ttc" experiment optkit-rag campaign serve \
  --store "$STORE" --listen "127.0.0.1:$PORT" \
  --workbench-token "$TOKEN" --workbench-actor "$ACTOR" \
  >"$TMP/server.log" 2>&1 &
SERVER_PID=$!
for _ in $(seq 1 100); do
  if curl -fsS "$BASE/api/rag/v1/health" >"$TMP/health.json" 2>/dev/null; then
    break
  fi
  sleep 0.05
done
curl -fsS "$BASE/api/rag/v1/health" >/dev/null

unauthorized=$(curl -sS -o "$TMP/unauthorized.json" -w '%{http_code}' "$BASE/api/rag/workbench/v1/catalog")
[[ "$unauthorized" == "401" ]]
curl -fsS -H "Authorization: Bearer $TOKEN" "$BASE/api/rag/workbench/v1/catalog" >"$TMP/catalog.json"

python3 - "$CAMPAIGN" "$TMP/compile.json" <<'PY'
import json,sys
campaign,out=sys.argv[1:]
body={"parent":{"campaign":campaign,"arm":"limit-1"},"mutations":[{"variable":"fusion.rrf_k","value":20}]}
open(out,'w').write(json.dumps(body))
PY
curl -fsS -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  --data-binary @"$TMP/compile.json" "$BASE/api/rag/workbench/v1/proposals:compile" >"$TMP/draft.json"
DRAFT=$(python3 -c 'import json,sys; v=json.load(sys.stdin); assert v["draft"]["sealable"]; print(v["draft"]["digest"])' <"$TMP/draft.json")

python3 - "$CAMPAIGN" "$DRAFT" "$TMP/preview.json" "$TMP/seal.json" <<'PY'
import json,sys
campaign,digest,preview_out,seal_out=sys.argv[1:]
parent={"campaign":campaign,"arm":"limit-1"}
mutations=[{"variable":"fusion.rrf_k","value":20}]
preview={"parent":parent,"mutations":mutations,"draft_digest":digest,"probe":"fusion.rrf-contributions/v1","case_id":"q-comparison"}
intent={"proposer":{"kind":"human"},"strategy":"manual-coordinate/v1","hypothesis":"Lower k should increase early-rank influence.","expected_improvement":{"metric":"retrieval.target-coverage","groups":["multi-source"]},"risks":["a noisy early rank may receive more influence"],"motivation":{"case_ids":["q-comparison"]}}
seal={"campaign":campaign,"arm_id":"rrf-20","parent_arm_id":"limit-1","parent":parent,"draft_digest":digest,"mutations":mutations,"intent":intent}
open(preview_out,'w').write(json.dumps(preview)); open(seal_out,'w').write(json.dumps(seal))
PY
curl -fsS -H "Authorization: Bearer $TOKEN" -H 'Content-Type: application/json' \
  --data-binary @"$TMP/preview.json" "$BASE/api/rag/workbench/v1/previews:run" >"$TMP/preview-response.json"
curl -fsS -H "Authorization: Bearer $TOKEN" -H 'Idempotency-Key: live-seal-1' -H 'Content-Type: application/json' \
  --data-binary @"$TMP/seal.json" "$BASE/api/rag/workbench/v1/proposals:seal" >"$TMP/sealed.json"
curl -fsS -H "Authorization: Bearer $TOKEN" -H 'Idempotency-Key: live-seal-1' -H 'Content-Type: application/json' \
  --data-binary @"$TMP/seal.json" "$BASE/api/rag/workbench/v1/proposals:seal" >"$TMP/retried.json"
cmp "$TMP/sealed.json" "$TMP/retried.json"
curl -fsS "$BASE/api/rag/v1/campaigns/$CAMPAIGN/cockpit" >"$TMP/cockpit.json"

python3 - "$TMP" "$ACTOR" "$CAMPAIGN" <<'PY'
import json,pathlib,sys
root=pathlib.Path(sys.argv[1]); actor=sys.argv[2]; campaign=sys.argv[3]
health=json.loads((root/'health.json').read_text())
catalog=json.loads((root/'catalog.json').read_text())
draft=json.loads((root/'draft.json').read_text())['draft']
preview=json.loads((root/'preview-response.json').read_text())['preview']
sealed=json.loads((root/'sealed.json').read_text())['sealed']
cockpit=json.loads((root/'cockpit.json').read_text())
error=json.loads((root/'unauthorized.json').read_text())['error']
assert health['read_only'] is True
assert catalog['catalog']['semantic_id']=='sha256:d3034d1d61cb5da92649bf9d199015e25a6e5223ed50741f594ceba2093730b6'
assert draft['sealable'] is True
assert preview['probe']=='fusion.rrf-contributions/v1' and preview['sensitivity']=='internal'
assert sealed['candidate']['proposer']['identity']==actor
assert cockpit['campaign']==campaign
assert error['code']=='unauthenticated'
print(f"campaign: {campaign}")
print("specialist_read_only: true")
print(f"catalog_semantic_id: {catalog['catalog']['semantic_id']}")
print(f"draft_digest: {draft['digest']}")
print(f"preview_probe: {preview['probe']}")
print(f"candidate_id: {sealed['candidate']['id']}")
print(f"candidate_actor: {sealed['candidate']['proposer']['identity']}")
print("unauthorized_catalog_status: 401")
print("idempotent_retry_equal: true")
print("result: PASS - live composed server read and command workflow")
PY
if [[ -n "${OUTPUT_DIR:-}" ]]; then
  mkdir -p "$OUTPUT_DIR"
  cp "$TMP/health.json" "$OUTPUT_DIR/specialist-health.json"
  cp "$TMP/catalog.json" "$OUTPUT_DIR/catalog.json"
  cp "$TMP/draft.json" "$OUTPUT_DIR/compile.json"
  cp "$TMP/preview-response.json" "$OUTPUT_DIR/preview.json"
  cp "$TMP/sealed.json" "$OUTPUT_DIR/seal.json"
  cp "$TMP/cockpit.json" "$OUTPUT_DIR/specialist-cockpit.json"
  cp "$TMP/unauthorized.json" "$OUTPUT_DIR/unauthorized.json"
fi
