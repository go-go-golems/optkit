#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
RAG_TTC="$WORKSPACE/rag-ttc"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT
STORE="$TMP/store"
MANIFEST="$TMP/ephemeral-manifest.yaml"
cp "$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v3.yaml" "$MANIFEST"

(
  cd "$RAG_TTC"
  GOWORK=off go build -o "$TMP/rag-ttc" ./cmd/rag-ttc
)
"$TMP/rag-ttc" experiment optkit-rag campaign run \
  --manifest "$MANIFEST" --store "$STORE" --reset --format json >"$TMP/run.json"
CAMPAIGN=$(python3 -c 'import json,sys; print(json.load(sys.stdin)[0]["campaign"])' <"$TMP/run.json")
rm "$MANIFEST"
[[ ! -e "$MANIFEST" ]]
"$TMP/rag-ttc" experiment optkit-rag campaign status \
  --campaign "$CAMPAIGN" --store "$STORE" --format json >"$TMP/status.json"
"$TMP/rag-ttc" experiment optkit-rag campaign verify \
  --campaign "$CAMPAIGN" --store "$STORE" --format json >"$TMP/verify.json"
"$TMP/rag-ttc" experiment optkit-rag campaign resume \
  --campaign "$CAMPAIGN" --store "$STORE" --format json >"$TMP/resume.json"

python3 - "$TMP" "$STORE/optkit.db" "$CAMPAIGN" <<'PY'
import json, pathlib, sqlite3, sys
root = pathlib.Path(sys.argv[1])
db_path, campaign = sys.argv[2], sys.argv[3]
status = json.loads((root / "status.json").read_text())[0]
verify = json.loads((root / "verify.json").read_text())[0]
resume = json.loads((root / "resume.json").read_text())[0]
assert status["campaign"] == resume["campaign"] == verify["campaign"] == campaign
assert status["status"] == resume["status"] == "completed"
assert status["events"] == resume["events"]
assert verify["journal_verified"] and verify["direct_payloads_verified"] and verify["nested_payloads_verified"]
conn = sqlite3.connect(db_path)
rows = conn.execute("SELECT kind,subject,payload_schema FROM campaign_events WHERE campaign_id=? AND kind IN ('CandidateProposed','SnapshotMaterialized') ORDER BY seq", (campaign,)).fetchall()
assert [row[0] for row in rows] == ["CandidateProposed", "SnapshotMaterialized"]
assert rows[0][2] == "schema:rag-ttc.candidate-proposal/v1"
assert rows[1][2] == "schema:optkit.snapshot/v1"
print(f"campaign: {campaign}")
print("manifest_present_after_run: false")
print(f"status_after_reopen: {status['status']}")
print(f"resume_status: {resume['status']}")
print(f"events_before_after_resume: {status['events']}/{resume['events']}")
print(f"candidate_id: {rows[0][1]}")
print(f"child_snapshot: {rows[1][1]}")
print(f"nested_payloads_verified: {verify['nested_payloads_verified']}")
print("result: PASS - source manifest removal did not affect restart or explanation facts")
PY
