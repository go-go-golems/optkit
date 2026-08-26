---
Title: Implementation Diary
Ticket: OPTKIT-007
Status: active
Topics: [implementation, optkit, rag-ttc]
DocType: reference
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: "Chronological implementation record for the specialist UI backend."
LastUpdated: 2026-08-26T00:30:00Z
WhatFor: "Preserve prompts, decisions, failures, commits, validation, and frontend handoff evidence."
WhenToUse: "Read before resuming or reviewing OPTKIT-007."
---

# Diary

## Goal

Build the verified, bounded, GET-only RAG-TTC backend that carries durable experiment facts up to the frontend integration boundary.

## Step 1: Freeze the read model and ticket plan

This step created the ticket, split implementation into six explicit phases, mapped existing HTTP composition roots and durable artifact schemas, and froze the first version of the specialist DTOs before projector implementation.

The API remains product-owned and read-only. Its first vertical path is cockpit → arm comparison → paired case → episode pipeline → episode provenance.

### Prompt Context

**User prompt (verbatim):**

```text
Ok, can we now Let's do the server backend side of things, all the way up to the UI, which we'll give to our frontend engineer.

commit at appropriate intervals and keep a detailed diary as you work (using the diary format from the skill)
Print out a brutalist work slip with the plan / different phases for the ticket. then before stsarting a phase, plrint a split about the phase, and print one when the phase is done.
```

**Assistant interpretation:** Implement the complete specialist read backend, characterize it end to end, and leave a precise frontend integration package, with phase slips, commits, and diary updates.

**Inferred user intent:** Give the frontend engineer stable, trustworthy APIs backed by the same durable facts already validated through the CLI.

### What I did

- Created OPTKIT-007 and six phase tasks.
- Printed the overall plan and Phase 0 start slips.
- Mapped RAG-TTC server roots, `http.ServeMux` conventions, Optkit query limits, campaign events, episode results, trajectories, and artifact sensitivity.
- Added versioned cockpit, comparison, case-page, pipeline, provenance, diagnostic, stage, and artifact DTOs.
- Added opaque cursor and bounded case-limit contracts.
- Wrote the initial API/projector design.

### Why

- DTOs and failure semantics must stabilize before multiple projector and HTTP implementations depend on them.
- The browser must not become a journal reducer or trajectory decoder.
- Explicit preview policy is necessary before returning artifact content.

### What worked

- Existing Optkit and RAG-TTC records contain enough information for a retrieval-focused vertical path.
- Go 1.22 method-aware `ServeMux` patterns match repository conventions.
- OPTKIT-006 already persists per-arm graphs and manifest identity in `CampaignSpec`.

### What didn't work

- N/A during initial contract creation.

### What I learned

- The current durable fixture can fully support retrieval pipeline and provenance, but not context/answer/judge stages.
- Direct event payloads are sufficient to find campaign spec, completions, observations, and estimates.
- Pipeline stage payloads are separately referenced from the sealed trajectory and can follow bounded preview policy.

### What was tricky to build

- Missing observations must have explicit `unknown` status and no numeric field. Zero is a valid measured result and cannot represent absence.
- Case pagination needs an opaque token even though the first fixture is small; this prevents the frontend from depending on array offsets.

### What warrants a second pair of eyes

- Review DTO field names before a frontend ships against v1.
- Review whether 4096 bytes is the right initial preview ceiling.
- Confirm first provenance kind should remain episode-only.

### What should be done in the future

- Add context, answer, and judge DTO fields only when durable producers emit them.
- Replace in-memory case slicing with database pagination for large campaigns.

### Code review instructions

- Start with `specialistapi/types.go` and the route list in the design doc.
- Confirm every top-level DTO carries API version, schema, and journal sequence.
- Run `go test ./pkg/ttc/specialistapi -count=1`.

### Technical details

```text
DefaultCaseLimit=50
MaximumCaseLimit=100
MaximumArtifactPreviewSize=4096
API version=rag-ttc.specialist-api/v1
```
