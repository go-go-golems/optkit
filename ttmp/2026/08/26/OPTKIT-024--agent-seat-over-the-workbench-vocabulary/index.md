---
Title: Agent Seat over the Workbench Vocabulary
Ticket: OPTKIT-024
Status: complete
Topics:
    - design
    - ui
    - rag-ttc
DocType: index
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Give an agent the same workbench vocabulary humans use — mentions as live presentations, agent draft verbs in the shared proposal document under actor attribution, sealing behind human approval.
LastUpdated: 2026-08-26T16:04:32.139075399-04:00
WhatFor: Extend "one vocabulary for playing and proposing" to non-human proposers without a parallel mutation mechanism.
WhenToUse: Read when implementing agent integration; strictly gated on OPTKIT-023 proving the human path.
---

# Agent Seat over the Workbench Vocabulary

## Overview

Final ticket of the PBUI track. Exports the OPTKIT-021 vocabulary in the
pbui-chat shape, renders agent mentions as live presentations, routes agent
draft verbs into the shared `ragttc.proposal-draft/v1` document with actor
attribution, and gates `proposal.seal` / `trial.run` behind human approval —
making "the agent cannot seal alone" an authorization fact. The sealed
candidate's Proposer records the agent; the verb trace records who did what.

**Program position:** gated on OPTKIT-023; deliberately last, because every
agent affordance is a re-use of the human path.

- [Intern guide](design-doc/01-intern-guide-agent-seat-over-the-workbench-vocabulary.md)

## Key Links

- **Related Files**: See frontmatter RelatedFiles field
- **External Sources**: See frontmatter ExternalSources field

## Status

Current status: **active**

## Topics

- design
- ui
- rag-ttc

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
