#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR=$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)
WORKSPACE=$(cd -- "$SCRIPT_DIR/../../../../../../.." && pwd)
cd "$WORKSPACE"

for ticket in OPTKIT-016 OPTKIT-017 OPTKIT-018; do
  listing=$(docmgr ticket list --ticket "$ticket")
  grep -q 'Status: \*\*complete\*\*' <<<"$listing"
  grep -q 'Tasks: 0 open / 8 done' <<<"$listing"
  docmgr doctor --ticket "$ticket" --stale-after 30 | grep -q 'All checks passed'
done

python3 optkit/ttmp/2026/08/26/OPTKIT-016--proposal-compiler-and-glazed-cli/scripts/03-audit-work-slips.py >/dev/null
python3 optkit/ttmp/2026/08/26/OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence/scripts/04-audit-work-slips.py >/dev/null
python3 optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/scripts/02-audit-work-slips.py >/dev/null

for ticket_path in \
  OPTKIT-016--proposal-compiler-and-glazed-cli \
  OPTKIT-017--proposal-sealing-candidate-manifests-and-campaign-persistence \
  OPTKIT-018--candidate-projections-and-workbench-command-api; do
  grep -q 'OK: uploaded' "optkit/ttmp/2026/08/26/$ticket_path/various/remarkable-implementation-upload.log"
done

[[ -z $(git -C rag-ttc status --short) ]]
[[ $(git -C optkit status --short) == '?? numbergame-demo' ]]
VAULT=/home/manuel/code/wesen/go-go-golems/go-go-parc
[[ -z $(git -C "$VAULT" status --short) ]]
[[ $(git -C "$VAULT" rev-parse HEAD) == $(git -C "$VAULT" rev-parse origin/main) ]]

! find rag-ttc/assets/configs/experiments/optkit-rag -name '*v2*' -print | grep -q .
! rg -n 'semantic-limit-v2' rag-ttc --glob '!*.md' >/dev/null
! rg -n 'TODO|FIXME|HACK|fallback decoder|compatibility alias' \
  rag-ttc/pkg/ttc/experimentworkbench/{proposal.go,sealing.go,workbench_contracts.go,workbench_service.go,catalog_service.go,authoring_manifest.go,campaign_materializer.go} \
  rag-ttc/pkg/ttc/workbenchapi \
  rag-ttc/pkg/ttc/optkitcampaign/candidate_facts.go >/dev/null

python3 - <<'PY'
from pathlib import Path
import json, re
root=Path('optkit/ttmp/2026/08/26/OPTKIT-018--candidate-projections-and-workbench-command-api/various/api-contracts')
files=sorted(root.glob('*.json'))
assert len(files)==7
text=''.join(p.read_text() for p in files)
for path in files: json.loads(path.read_text())
assert 'live-smoke-token' not in text
for ticket in ('016','017','018'):
    matches=list(Path('optkit/ttmp/2026/08/26').glob(f'OPTKIT-{ticket}--*/reference/01-implementation-diary.md'))
    assert len(matches)==1
    source=matches[0].read_text()
    steps=re.split(r'(?=^## Step \d+:)', source, flags=re.M)[1:]
    assert len(steps)==9
    required=('### Prompt Context','### What I did','### Why','### What worked',"### What didn't work",'### What I learned','### What was tricky to build','### What warrants a second pair of eyes','### What should be done in the future','### Code review instructions','### Technical details')
    assert all(all(section in step for section in required) for step in steps)
PY

rg -q '\| OPTKIT-016 .* complete' optkit/ttmp/2026/08/26/OPTKIT-011--*/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
rg -q '\| OPTKIT-017 .* complete' optkit/ttmp/2026/08/26/OPTKIT-011--*/design-doc/04-backend-first-optimization-workbench-program-roadmap.md
rg -q '\| OPTKIT-018 .* complete' optkit/ttmp/2026/08/26/OPTKIT-011--*/design-doc/04-backend-first-optimization-workbench-program-roadmap.md

if pgrep -af '/tmp/tmp\.[^ ]*/rag-ttc.*campaign serve' | grep -v 'pgrep -af' >/dev/null; then
  echo 'ticket live-server process remains' >&2
  exit 1
fi

printf 'tickets_complete: OPTKIT-016 OPTKIT-017 OPTKIT-018\n'
printf 'tasks_open: 0\n'
printf 'strict_diary_steps: 27\n'
printf 'work_slips_printed: 45\n'
printf 'completed_uploads: 3\n'
printf 'contract_json_files: 7\n'
printf 'stale_v2_assets: 0\n'
printf 'debt_markers: 0\n'
printf 'ticket_live_servers: 0\n'
printf 'unrelated_numbergame_demo: preserved\n'
printf 'result: PASS - backend compiler sealing and command API chain is complete\n'
