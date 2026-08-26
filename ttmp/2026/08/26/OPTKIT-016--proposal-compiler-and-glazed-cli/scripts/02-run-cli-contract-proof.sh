#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
RAG_TTC="$WORKSPACE/rag-ttc"
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

(
  cd "$RAG_TTC"
  GOWORK=off go build -o "$TMP/rag-ttc" ./cmd/rag-ttc
)

MANIFEST="$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v2.yaml"
"$TMP/rag-ttc" experiment optkit-rag catalog list --format json >"$TMP/catalog.json"
"$TMP/rag-ttc" experiment optkit-rag catalog show --variable fusion.rrf_k --format json >"$TMP/variable.json"
"$TMP/rag-ttc" experiment optkit-rag proposal compile \
  --manifest "$MANIFEST" --parent limit-2 \
  --set fusion.rrf_k=20 --format json >"$TMP/proposal.json"
"$TMP/rag-ttc" experiment optkit-rag proposal compile \
  --manifest "$MANIFEST" --parent limit-2 \
  --set fusion.rrf_k=oops --format json >"$TMP/invalid.json"

python3 - "$TMP" <<'PY'
import json, pathlib, sys
root = pathlib.Path(sys.argv[1])
catalog = json.loads((root / "catalog.json").read_text())
variable = json.loads((root / "variable.json").read_text())
proposal = json.loads((root / "proposal.json").read_text())
invalid = json.loads((root / "invalid.json").read_text())
assert [row["variable"] for row in catalog] == ["retrieval.final_result_limit", "fusion.rrf_k"]
assert len({row["semantic_catalog_id"] for row in catalog}) == 1
assert variable[0]["descriptor"]["value"]["float_range"] == {"minimum": 0.001, "maximum": 1000}
assert variable[0]["default"] == 60
assert len(proposal) == 1 and proposal[0]["sealable"] is True
assert proposal[0]["child_config"]["fusion"]["rrf_k"] == 20
assert proposal[0]["invalidation_plan"]["changes"] == ["fusion"]
assert proposal[0]["preview_capabilities"] == [{
    "variable": "fusion.rrf_k", "mode": "deterministic_local", "probe": "fusion.rrf-contributions/v1"
}]
assert len(invalid) == 1 and invalid[0]["sealable"] is False
assert [entry["code"] for entry in invalid[0]["diagnostics"]] == ["invalid_json"]
print(f"catalog_rows: {len(catalog)}")
print(f"semantic_catalog_id: {catalog[0]['semantic_catalog_id']}")
print(f"full_catalog_id: {catalog[0]['catalog_id']}")
print(f"rrf_domain: {variable[0]['descriptor']['value']['float_range']}")
print(f"draft_digest: {proposal[0]['digest']}")
print(f"before_graph: {proposal[0]['before_graph']['id']}")
print(f"after_graph: {proposal[0]['after_graph']['id']}")
print(f"direct_changes: {proposal[0]['invalidation_plan']['changes']}")
print(f"invalid_diagnostic: {invalid[0]['diagnostics'][0]['code']}")
print("result: PASS - catalog and proposal CLI preserve structured contracts")
PY
