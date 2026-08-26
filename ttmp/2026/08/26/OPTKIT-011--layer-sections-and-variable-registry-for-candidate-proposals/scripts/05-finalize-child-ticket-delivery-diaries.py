#!/usr/bin/env python3
"""Append validation, commit, and reMarkable delivery evidence to child diaries."""
from pathlib import Path

ROOT = Path("optkit/ttmp/2026/08/26")
DATA = {
    "OPTKIT-012": ("architecture-closure-and-optimization-workbench-contracts", "428f6b8f5391dc851d989364b21a9d727c78cfc0", "OPTKIT-012-014: design workbench foundations", "OPTKIT-012 Architecture Contracts Guide"),
    "OPTKIT-013": ("optkit-semantic-catalog-and-executable-variable-bindings", "428f6b8f5391dc851d989364b21a9d727c78cfc0", "OPTKIT-012-014: design workbench foundations", "OPTKIT-013 Catalog and Bindings Guide"),
    "OPTKIT-014": ("whole-pipeline-rag-configuration-and-graph-derivation", "428f6b8f5391dc851d989364b21a9d727c78cfc0", "OPTKIT-012-014: design workbench foundations", "OPTKIT-014 Pipeline Config Guide"),
    "OPTKIT-015": ("real-fusion-configuration-and-first-rag-optimization-catalog", "83d0f4f201ae59a0d8983d476e7873312fca8142", "OPTKIT-015-017: design RAG proposal backend", "OPTKIT-015 Fusion and RAG Catalog Guide"),
    "OPTKIT-016": ("proposal-compiler-and-glazed-cli", "83d0f4f201ae59a0d8983d476e7873312fca8142", "OPTKIT-015-017: design RAG proposal backend", "OPTKIT-016 Proposal Compiler CLI Guide"),
    "OPTKIT-017": ("proposal-sealing-candidate-manifests-and-campaign-persistence", "83d0f4f201ae59a0d8983d476e7873312fca8142", "OPTKIT-015-017: design RAG proposal backend", "OPTKIT-017 Proposal Sealing Guide"),
    "OPTKIT-018": ("candidate-projections-and-workbench-command-api", "ec784b0d84e47ad45a7429f98482e9678df5dac7", "OPTKIT-018-020: design workbench delivery layers", "OPTKIT-018 Workbench API Guide"),
    "OPTKIT-019": ("react-workbench-framework-and-rrf-vertical-slice", "ec784b0d84e47ad45a7429f98482e9678df5dac7", "OPTKIT-018-020: design workbench delivery layers", "OPTKIT-019 React Workbench RRF Guide"),
    "OPTKIT-020": ("asset-variable-proof-for-representation-prompts", "ec784b0d84e47ad45a7429f98482e9678df5dac7", "OPTKIT-018-020: design workbench delivery layers", "OPTKIT-020 Prompt Asset Variable Guide"),
}

for ticket, (slug, commit, message, bundle) in DATA.items():
    directory = ROOT / f"{ticket}--{slug}"
    diary = directory / "reference/01-implementation-diary.md"
    text = diary.read_text()
    if "## Step 2: Validate, commit, and deliver the guide" in text:
        print(f"skip existing {ticket}")
        continue
    addition = f'''

## Step 2: Validate, commit, and deliver the guide

This step converted the researched guide from a working document into a reviewed ticket deliverable. The ticket's frontmatter, relations, tasks, and changelog were validated; the documentation was committed in a dependency-coherent batch; and the index, guide, and diary were rendered and uploaded as one reMarkable PDF with a table of contents.

The implementation tasks intentionally remain open. This delivery completes the up-front planning package and gives the future implementer an evidence-backed starting point, not a false claim that production behavior has already changed.

### Prompt Context

**User prompt (verbatim):** (see Step 1)

**Assistant interpretation:** Finish the ticket documentation package with validation, coherent Git history, strict diary evidence, and a ticket-specific reMarkable delivery.

**Inferred user intent:** Make the plan durable and reviewable before implementation starts, while preserving a clear distinction between completed design work and pending code tasks.

**Commit (documentation):** `{commit}` — "{message}"

### What I did

- Ran `docmgr doctor --ticket {ticket} --stale-after 30` and obtained `All checks passed`.
- Scanned the index, guide, and diary for generated placeholder sections; none remained.
- Ran focused baseline Go tests before documentation changes; all selected Optkit and RAG-TTC packages passed.
- Ran `git diff --check`/`git diff --cached --check`, corrected whitespace findings, and committed the ticket package.
- Ran the program upload script in `--dry-run` mode, then rendered and uploaded `{bundle}.pdf`.
- Preserved upload evidence in `various/remarkable-dry-run.log` and `various/remarkable-upload.log`.

### Why

- A detailed guide is only useful when frontmatter, links, task state, and delivery artifacts agree.
- Grouping commits by dependency layer keeps review focused while avoiding one nine-ticket mega-commit.
- The reMarkable bundle lets the architecture be reviewed away from the source tree without losing the index or diary context.

### What worked

- `remarquee` reported: `OK: uploaded {bundle}.pdf -> /ai/2026/08/26/{ticket}`.
- The upload rendered the real Markdown rather than only performing a path dry-run.
- The ticket retains open implementation tasks and checked documentation/delivery tasks separately.

### What didn't work

- The first staged diff check failed with `new blank line at EOF` in generated changelogs and `trailing whitespace` on blank quoted prompt lines (`+> `) in generated diaries. The diary generator was corrected to emit `>` on blank quote lines, changelog EOFs were normalized, and the second `git diff --check` passed.
- OPTKIT-020 initially referenced nonexistent `repo://optkit/artifact/artifact.go`; the actual package contains `ref.go`, `store.go`, `read.go`, and `helpers.go`. The invalid relation was removed, `artifact/ref.go` was related, and doctor then passed.

### What I learned

- A dry-run validates upload selection and destination but the real upload is the evidence that Pandoc/LaTeX can render the complete guide.
- Generated prose containing blockquotes needs whitespace validation just like source code.
- Keeping implementation tasks open while checking documentation/delivery tasks makes ticket status truthful.

### What was tricky to build

- The uploaded diary necessarily describes the work up to its render time. This final local step records the upload receipt after the PDF has been created; re-uploading with `--force` solely to include its own receipt would overwrite a new document and risk future annotations.
- Cross-ticket commit hashes and ticket-specific upload names had to remain consistent across nine independent workspaces.

### What warrants a second pair of eyes

- Review the proposed APIs and compatibility decisions before implementation; successful document delivery is not architecture acceptance.
- Confirm the reMarkable bundle name and ticket folder are the intended long-term review locations before adding annotations.

### What should be done in the future

- Begin only after upstream entry gates in the OPTKIT-011 roadmap are satisfied.
- During implementation, append new diary steps with exact code commit hashes, failures, commands, and fresh validation evidence.

### Code review instructions

- Start with `index.md`, then read the complete `design-doc/01-*.md`, then this diary.
- Inspect commit `{commit}` for the documentation batch.
- Validate locally with `docmgr doctor --ticket {ticket} --stale-after 30`.
- Consult `tasks.md` for the still-open implementation sequence.

### Technical details

```text
bundle: {bundle}.pdf
remote: /ai/2026/08/26/{ticket}
commit: {commit}
doctor: clean
production code changes: none
```
'''
    diary.write_text(text.rstrip() + addition.rstrip() + "\n")
    print(diary)
