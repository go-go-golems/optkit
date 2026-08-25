#!/usr/bin/env python3
from pathlib import Path
import sys
import yaml

script = Path(__file__).resolve()
optkit = script.parents[6]
doc = optkit / "docs/01-optkit-records-artifacts-and-control-model.md"
text = doc.read_text()

if not text.startswith("---\n"):
    raise SystemExit("document does not start with YAML frontmatter")
_, frontmatter, body = text.split("---\n", 2)
data = yaml.safe_load(frontmatter)
required = {
    "Title", "Slug", "Short", "Topics", "Commands", "Flags",
    "IsTopLevel", "IsTemplate", "ShowPerDefault", "SectionType",
}
missing = sorted(required - set(data))
if missing:
    raise SystemExit(f"missing Glazed frontmatter fields: {missing}")
if data["SectionType"] != "GeneralTopic":
    raise SystemExit("field reference must be a GeneralTopic")
if any(line.startswith("# ") for line in body.splitlines()):
    raise SystemExit("Glazed help entry must not contain a top-level Markdown heading")
for heading in ("## Troubleshooting", "## See Also", "## Completed removal register"):
    if heading not in body:
        raise SystemExit(f"missing required section: {heading}")

required_terms = [
    "record.Digest", "artifact.Ref", "VariableDescriptor", "SnapshotRecord",
    "Candidate", "system.Factory", "episode.Event", "EpochDefinition",
    "Observation", "DatasetManifest", "TrialPlan", "EpisodeSpec", "Estimate",
    "campaign.Command", "ControlEvent", "campaign.State", "Journal",
    "WorkItem", "LeaseRequest", "WorkRecord", "ResourceClaim", "Heartbeat",
    "Reservation", "projection.Overview", "CampaignSummary", "EventView",
    "local.Profile", "store/sqlite.Store",
]
missing_terms = [term for term in required_terms if term not in text]
if missing_terms:
    raise SystemExit(f"missing model terms: {missing_terms}")

legacy = [
    optkit / "docs/implementation-diary.md",
    optkit / "docs/adr",
    optkit / "docs/journals",
]
remaining = [str(path) for path in legacy if path.exists()]
if remaining:
    raise SystemExit(f"legacy implementation history remains in docs: {remaining}")

slug = data["Slug"]
for other in (optkit / "docs").rglob("*.md"):
    if other == doc:
        continue
    if f"Slug: {slug}" in other.read_text():
        raise SystemExit(f"duplicate Glazed slug in {other}")

print("GLAZED_FRONTMATTER=PASS")
print(f"SLUG={slug}")
print(f"LINES={len(text.splitlines())}")
print(f"WORDS={len(text.split())}")
print("LEGACY_DOC_RELOCATION=PASS")
print("MODEL_DOC_VALIDATION=PASS")
