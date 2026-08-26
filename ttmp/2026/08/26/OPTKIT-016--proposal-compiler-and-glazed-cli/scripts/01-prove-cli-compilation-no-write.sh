#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
RAG_TTC="$WORKSPACE/rag-ttc"
RUNS=${RUNS:-100}
TMP=$(mktemp -d)
trap 'rm -rf "$TMP"' EXIT

PROOF_ROOT="$TMP/watched"
mkdir -p "$PROOF_ROOT"
cp "$RAG_TTC/assets/configs/experiments/optkit-rag/semantic-limit-v3.yaml" "$PROOF_ROOT/manifest.yaml"
cat >"$PROOF_ROOT/mutations.yaml" <<'YAML'
- variable: fusion.rrf_k
  value: 20
YAML

(
  cd "$RAG_TTC"
  GOWORK=off go build -o "$TMP/rag-ttc" ./cmd/rag-ttc
)

tree_digest() {
  find "$PROOF_ROOT" -type f -print0 \
    | sort -z \
    | xargs -0 sha256sum \
    | sha256sum \
    | awk '{print $1}'
}

before=$(tree_digest)
expected=""
for ((run = 1; run <= RUNS; run++)); do
  output=$(
    cd "$PROOF_ROOT"
    "$TMP/rag-ttc" experiment optkit-rag proposal compile \
      --manifest manifest.yaml \
      --parent limit-2 \
      --mutations mutations.yaml \
      --format json
  )
  digest=$(python3 -c 'import json,sys; rows=json.load(sys.stdin); assert len(rows)==1 and rows[0]["sealable"] is True; print(rows[0]["digest"])' <<<"$output")
  if [[ -z "$expected" ]]; then
    expected=$digest
  elif [[ "$digest" != "$expected" ]]; then
    printf 'draft digest changed on run %d: %s != %s\n' "$run" "$digest" "$expected" >&2
    exit 1
  fi
done
after=$(tree_digest)

[[ "$before" == "$after" ]] || {
  printf 'watched filesystem changed: %s != %s\n' "$before" "$after" >&2
  find "$PROOF_ROOT" -maxdepth 4 -printf '%P %y %s\n' >&2
  exit 1
}
[[ $(find "$PROOF_ROOT" -type f | wc -l) -eq 2 ]]
[[ $(find "$PROOF_ROOT" \( -name '*.sqlite' -o -name '*.db' -o -name 'artifacts' -o -name 'journal' \) | wc -l) -eq 0 ]]

printf 'runs: %d\n' "$RUNS"
printf 'draft_digest: %s\n' "$expected"
printf 'before_tree_sha256: %s\n' "$before"
printf 'after_tree_sha256: %s\n' "$after"
printf 'watched_files: 2\n'
printf 'sqlite_files: 0\n'
printf 'artifact_directories: 0\n'
printf 'journal_paths: 0\n'
printf 'result: PASS - repeated CLI compilation performed no durable writes\n'
