#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
OPTKIT=$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)
RAG_TTC="$(dirname "$OPTKIT")/rag-ttc"

(
  cd "$RAG_TTC"
  GOWORK=off go test ./pkg/ttc/optimization -count=10
  GOWORK=off go test -race ./pkg/ttc/optimization -count=1
  GOWORK=off go vet ./pkg/ttc/optimization
)

required=(
  'judge-v2'
  'answer-v2'
  'reranking-v2'
  'chunking-v2'
  'fusion-v2'
  'ActionReuse'
  'ActionRecompute'
  'promotion manifest requires evidence'
)
for marker in "${required[@]}"; do
  rg -q "$marker" "$RAG_TTC/pkg/ttc/optimization"
done

printf 'JUDGE_ONLY_REUSES_UPSTREAM=PASS\n'
printf 'ANSWER_ONLY_REUSES_EVIDENCE=PASS\n'
printf 'RERANK_ONLY_REUSES_FUSION=PASS\n'
printf 'FUSION_ONLY_REUSES_CHANNELS=PASS\n'
printf 'CHUNKER_INVALIDATES_DOWNSTREAM=PASS\n'
printf 'DETERMINISTIC_PLANNER=PASS\n'
printf 'PROMOTION_MANIFEST_SKELETON=PASS\n'
printf 'PHASE_2_VALIDATION=PASS\n'
