#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
RAG_TTC="$WORKSPACE/rag-ttc"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
STORE="$TMP/store"
MANIFEST="$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v3.yaml"

(
  cd "$RAG_TTC"
  GOWORK=off go build -o "$TMP/rag-ttc" ./cmd/rag-ttc
)

"$TMP/rag-ttc" experiment optkit-rag campaign dry-run \
  --manifest "$MANIFEST" --store "$STORE" --format json >"$TMP/dry-run.json"
[[ ! -e "$STORE" ]]
"$TMP/rag-ttc" experiment optkit-rag campaign run \
  --manifest "$MANIFEST" --store "$STORE" --reset --format json >"$TMP/run.json"
CAMPAIGN=$(python3 -c 'import json,sys; rows=json.load(sys.stdin); assert len(rows)==1 and rows[0]["status"]=="completed"; print(rows[0]["campaign"])' <"$TMP/run.json")
"$TMP/rag-ttc" experiment optkit-rag campaign status \
  --campaign "$CAMPAIGN" --store "$STORE" --format json >"$TMP/status.json"
"$TMP/rag-ttc" experiment optkit-rag campaign verify \
  --campaign "$CAMPAIGN" --store "$STORE" --format json >"$TMP/verify.json"

python3 - "$TMP" "$STORE/optkit.db" "$CAMPAIGN" <<'PY'
import json, pathlib, sqlite3, sys
root = pathlib.Path(sys.argv[1])
db_path, campaign = sys.argv[2], sys.argv[3]
dry = json.loads((root / "dry-run.json").read_text())[0]
run = json.loads((root / "run.json").read_text())[0]
status = json.loads((root / "status.json").read_text())[0]
verify = json.loads((root / "verify.json").read_text())[0]
assert dry["mutation"] is False and dry["arms"] == 2 and dry["episodes"] == 6
assert run["campaign"] == status["campaign"] == verify["campaign"] == campaign
assert status["status"] == "completed" and status["completed"] == 6
assert verify["journal_verified"] is True
assert verify["direct_payloads_verified"] is True
assert verify["nested_payloads_verified"] is True
assert verify["unique_nested_payloads"] > 0
conn = sqlite3.connect(db_path)
counts = dict(conn.execute("SELECT kind, COUNT(*) FROM campaign_events WHERE campaign_id=? GROUP BY kind", (campaign,)))
assert counts["CandidateProposed"] == 1
assert counts["SnapshotMaterialized"] == 1
candidate = conn.execute("SELECT command_id, subject, payload_schema FROM campaign_events WHERE campaign_id=? AND kind='CandidateProposed'", (campaign,)).fetchone()
snapshot = conn.execute("SELECT command_id, subject, payload_schema FROM campaign_events WHERE campaign_id=? AND kind='SnapshotMaterialized'", (campaign,)).fetchone()
assert candidate[0] == snapshot[0] and candidate[0]
assert candidate[2] == "schema:rag-ttc.candidate-proposal/v1"
assert snapshot[2] == "schema:optkit.snapshot/v1"
command = conn.execute("SELECT first_seq,last_seq FROM campaign_commands WHERE campaign_id=? AND command_id=?", (campaign,candidate[0])).fetchone()
assert command[1] - command[0] == 1
print(f"campaign: {campaign}")
print(f"status: {status['status']}")
print(f"events: {verify['events']}")
print(f"candidate_events: {counts['CandidateProposed']}")
print(f"snapshot_events: {counts['SnapshotMaterialized']}")
print(f"candidate_command: {candidate[0]}")
print(f"candidate_id: {candidate[1]}")
print(f"child_snapshot: {snapshot[1]}")
print(f"unique_direct_payloads: {verify['unique_direct_payloads']}")
print(f"unique_nested_payloads: {verify['unique_nested_payloads']}")
print("result: PASS - campaign CLI compiled then sealed one durable candidate")
PY
