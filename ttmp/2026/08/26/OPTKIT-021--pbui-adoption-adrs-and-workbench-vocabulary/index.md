---
Title: PBUI Adoption ADRs and Workbench Vocabulary
Ticket: OPTKIT-021
Status: active
Topics:
    - architecture
    - design
    - ui
    - rag-ttc
    - optkit
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Architecture-closure ticket for building the optimization workbench frontend on PBUI, deciding product packaging, Go workbench service placement, document formats, and the presentation-type/verb vocabulary before UI implementation.
LastUpdated: 2026-08-26T16:04:31.934141225-04:00
WhatFor: Play the OPTKIT-012 role for the frontend so the type map, verbs, and package boundaries do not emerge accidentally across PRs.
WhenToUse: Read before starting OPTKIT-022 or OPTKIT-023, and whenever a new presentation type, verb, or workbench document format is proposed.
---

# PBUI Adoption ADRs and Workbench Vocabulary

## Overview

The program replaced OPTKIT-019's bespoke React shell with the PBUI stack
(`@hyperslop-systems/pbui`, `pbui-workbench`, `workbench-protocol`, `plot`).
This ticket closes the frontend architecture decisions before code exists:

- **ADR G** — product packaging: new pbui product at `rag-ttc/apps/workbench/web`.
- **ADR H** — Go workbench document service: `rag-ttc/pkg/ttc/workbenchhost`,
  mounted beside (never inside) `specialistapi`.
- **ADR I** — workbench document formats: `ragttc.proposal-draft/v1`,
  `ragttc.comparison/v1`, `ragttc.focus/v1`, with the derived-state rule
  (compile output is never stored).
- **ADR J** — the presentation vocabulary: read-side and authoring types,
  tones, conversions, environment.
- **ADR K** — verbs and the verb sink: navigation/draft/command families,
  danger verbs, the pure-compile debounce.
- **ADR L** — agent vocabulary reservation for OPTKIT-024.

**Program position:** entry gate for the PBUI track. Successor of the
superseded OPTKIT-019 together with OPTKIT-022 and OPTKIT-023. Runs in
parallel with the backend chain OPTKIT-015–018.

- [Intern guide (ADRs and vocabulary)](design-doc/01-intern-guide-pbui-adoption-adrs-and-workbench-vocabulary.md)
- [Diary (whole PBUI-track working session)](reference/01-diary.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- architecture
- design
- ui
- rag-ttc
- optkit

## Tasks

See [tasks.md](./tasks.md) for the current task list.

## Changelog

See [changelog.md](./changelog.md) for recent changes and decisions.

## Structure

- design/ - Architecture and design documents
- reference/ - Prompt packs, API contracts, context summaries
- playbooks/ - Command sequences and test procedures
- scripts/ - Temporary code and tooling
- various/ - Working notes and research
- archive/ - Deprecated or reference-only artifacts
