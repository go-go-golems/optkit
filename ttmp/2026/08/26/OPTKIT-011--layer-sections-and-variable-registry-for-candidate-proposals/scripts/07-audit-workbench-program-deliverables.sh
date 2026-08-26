#!/usr/bin/env bash
set -euo pipefail

ROOT="optkit/ttmp/2026/08/26"
ids=(OPTKIT-011 OPTKIT-012 OPTKIT-013 OPTKIT-014 OPTKIT-015 OPTKIT-016 OPTKIT-017 OPTKIT-018 OPTKIT-019 OPTKIT-020)

echo "== docmgr doctor =="
for id in "${ids[@]}"; do
  docmgr doctor --ticket "$id" --stale-after 30
 done

echo "== child guide and diary structure =="
for id in OPTKIT-{012..020}; do
  dir=$(find "$ROOT" -maxdepth 1 -type d -name "$id--*" -print -quit)
  guide=$(find "$dir/design-doc" -maxdepth 1 -type f -name '*.md' -print -quit)
  diary=$(find "$dir/reference" -maxdepth 1 -type f -name '*.md' -print -quit)
  lines=$(wc -l < "$guide")
  if (( lines < 300 )); then echo "$id guide too short: $lines" >&2; exit 1; fi
  for pattern in "Executive summary" "Implementation" "Test" "Risk" "File reference"; do
    grep -qi "$pattern" "$guide" || { echo "$id guide missing $pattern" >&2; exit 1; }
  done
  for pattern in "## Step 1:" "## Step 2:" "### Prompt Context" "### What I did" "### What didn't work" "### What was tricky to build" "### What warrants a second pair of eyes" "### What should be done in the future" "### Code review instructions" "### Technical details"; do
    grep -q "$pattern" "$diary" || { echo "$id diary missing $pattern" >&2; exit 1; }
  done
  grep -q "OK: uploaded" "$dir/various/remarkable-upload.log" || { echo "$id missing upload receipt" >&2; exit 1; }
  echo "$id: guide=${lines} lines diary=$(wc -l < "$diary") lines upload=OK"
done

echo "== placeholder scan =="
if rg -n "Provide a brief|Describe the problem|Outline the steps|Link to related|<!--" \
    "$ROOT"/OPTKIT-0{12,13,14,15,16,17,18,19,20}--*/{index.md,design-doc,reference}; then
  echo "generated placeholder found" >&2
  exit 1
fi

echo "== upload receipts =="
grep -h "OK: uploaded" "$ROOT"/OPTKIT-0{12,13,14,15,16,17,18,19,20}--*/various/remarkable-upload.log
PARENT="$ROOT/OPTKIT-011--layer-sections-and-variable-registry-for-candidate-proposals"
grep "OK: uploaded" "$PARENT/various/remarkable-upload.log" "$PARENT/various/remarkable-program-diary-upload.log"

echo "== baseline focused tests =="
(
  cd optkit
  GOWORK=off go test ./space ./examples/numbergame -count=1
)
(
  cd rag-ttc
  GOWORK=off go test ./pkg/ttc/optimization ./pkg/ttc/experimentworkbench \
    ./pkg/ttc/optkitcampaign ./pkg/ttc/search ./pkg/ttc/specialistapi -count=1
)

echo "AUDIT OK: roadmap + nine child ticket packages + diaries + uploads + baseline tests"
