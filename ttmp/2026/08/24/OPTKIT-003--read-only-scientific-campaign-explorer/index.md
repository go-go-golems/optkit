---
Title: Read-Only Scientific Campaign Explorer
Ticket: OPTKIT-003
Status: complete
Topics:
    - optkit
    - ui
    - query-plane
    - visualization
    - scientific-workflow
    - architecture
    - local-development
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Architecture, implementation, validation, and onboarding package for the dependency-free read-only Optkit campaign explorer.
LastUpdated: 2026-08-24T20:28:09.384619363-04:00
WhatFor: Navigate the explorer guide, diary, implementation evidence, scripts, and delivery state.
WhenToUse: Start here before reviewing or extending the Optkit query API or scientific explorer.
---


# Read-Only Scientific Campaign Explorer

## Overview

This ticket explains and implements a first read-only scientific explorer for CLI-authored Optkit campaigns. The explorer uses a standard-library Go GET/SSE query server and directly embedded semantic HTML, pure-white monochrome CSS, and plain JavaScript. It has no browser mutation API, Node dependency, framework, simulated window chrome, menu bar, or frontend build step.

## Key links

- [Architecture and implementation guide](design-doc/01-scientific-campaign-explorer-architecture-and-implementation-guide.md)
- [Investigation and implementation diary](reference/01-investigation-diary.md)
- [Numbergame evidence capture](sources/01-numbergame-evidence.md)
- [Explorer validation transcript](sources/02-explorer-validation.txt)
- [Evidence capture script](scripts/01-capture-numbergame-evidence.sh)
- [Explorer validation script](scripts/02-validate-explorer.sh)
- [Tasks](tasks.md)
- [Changelog](changelog.md)

## Implemented v0

- Stable local campaign listing.
- Verified campaign overview and budget projection.
- Bounded event JSON with safe visible event-payload previews.
- GET-only API and sequence-cursor SSE.
- `optkit serve` command.
- Campaign list and hash navigation.
- Search/filter, stage rail, facts, lineage, trial matrix, paired analysis, budgets, event inventory, timeline, and evidence drawer.
- Pure white, flat retro-monochrome presentation with modern fonts and text-only color accents.
- Full, no-CGO, race, lint, HTTP, SSE, mutation-denial, and browser smoke validation.

## Current status

- Research and design: complete.
- Working v0 implementation: complete.
- Documentation and diary: complete.
- Docmgr validation: clean.
- reMarkable delivery: uploaded to `/ai/2026/08/25/OPTKIT-003`.
