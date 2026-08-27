---
Title: Specialist Work Scenarios and Questions
Ticket: OPTKIT-009
Status: active
Topics:
    - implementation
    - rag-ttc
    - ui
DocType: design-doc
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources:
    - "https://hyperslop.systems/ — accent palette source (green 2db878, purple 805bd7, red ef4038, yellow f2ad00)"
Summary: The specialist UI reframed around its human user — someone making a chatbot's retrieval better — with judge traces, verdicts, and failure examples at the center, data graphics that answer questions, and prose that explains what things are.
WhatFor: Drive the evolution of the specialist UI from journal-shaped tables toward a tool people do work in.
WhenToUse: Read when prioritizing specialist UI features or projector/backend asks.
---

# Specialist Work Scenarios and Questions

Everything here starts from one person: someone whose job is to make a
chatbot's answers better, which mostly means making its **retrieval** better.
The journal, the projectors, the digests — all of that scaffolding exists to
serve that person's day. The UI must not be shaped like the journal; it must
be shaped like the work.

The work, concretely, is a loop:

> Read failures → form a hypothesis → change one thing → re-run → see which
> failures got fixed and what broke → repeat.

Averages start that loop; **examples** drive it. Nobody improves a retriever
by staring at 0.833; they improve it by reading the question the bot fumbled,
the evidence it actually fetched, the evidence it should have fetched, and
what the judge said about the difference.

## Scenarios, in the user's voice

### A. "Show me what's broken"

I don't open the tool to admire a completed campaign. I open it to find the
cases where the bot let a user down. The first thing I want is a **failure
gallery**: the questions with the worst verdicts, each showing the question
text, what came back, what was expected, and the judge's words for why it
fell short. Sorted by how badly, and by how often (a case that fails on every
repeat is a bug; one that fails sometimes is flakiness).

### B. "Says who?" — reading judge traces

Behind every score is a judge's verdict. Before I trust 0.5 I want to read
the trace: what the judge was shown, what it concluded, in its own words. And
over time I want calibration — when the judge and I disagree on ten examples,
that's a judge problem, not a retrieval problem, and fixing retrieval against
a miscalibrated judge is wasted work.

### C. "Why did this one fail?" — the autopsy

For one failing question I want the full story in one view: the question, the
chunks each stage fetched/merged/filtered/ranked, where the right chunk died
(or was never fetched at all), and the chunk texts so I can judge relevance
myself. The chunk-presence grid and the narrowing-funnel line answer "where";
the chunk text answers "why". A failure autopsy is: *where* did the evidence
go, then *should* it have survived.

### D. "Did my fix work — and what did it break?"

After a change I don't want just a delta on the mean. I want **churn**:
which failing cases now pass (fixed), which passing cases now fail
(regressions), which still fail (untouched). A +0.167 mean that hides one
regression is a worse trade than it looks. The paired slope graph shows this
today; a fixed/broke/still-broken triage list would say it in words.

### E. "Let me try it myself"

Intuition comes from contact. I want to type my own question at a sealed
arm and watch the pipeline handle it live — not recorded, not scored, just
me poking the system the way a user would. Five hand-run queries teach me
more about a reranker than any table.

### F. "What exactly am I testing?" — in words

Arms, cases, and campaigns should describe themselves in prose, not slugs.
`limit-2` should say "same recipe, but retrieval keeps 2 results instead
of 1". `q-policy-negative` should say "asks for restricted content; passing
means retrieving nothing". These sentences belong in the manifests and case
definitions and should flow through the API — the person writing the
experiment knows the why; the tool should carry it to the person reading the
results.

### G. "Can I cite this?" — trust and bookkeeping

When results go into a write-up, I need the boring things to be solid:
journal verified, no hidden failures, budget honest, and every number
traceable to episodes and artifacts. Identifiers stay one disclosure away,
never in the reading line.

## The questions, ranked by how often they come up

Every session:

1. Which questions is the bot failing, and how badly?
2. What did the judge actually say about each failure?
3. For this failure — where in the pipeline did the right evidence die?
4. What does this chunk/case actually say? (text, not IDs)
5. Did my last change fix failures, and did it break anything that worked?

Regularly:

6. What am I actually testing, in words? (arm and case descriptions)
7. What did my config change cost — what recomputes, what's reusable?
8. Do I have budget left for another sweep?
9. Do I trust the judge? (calibration, spot-check disagreements)

Occasionally:

10. Let me hand-run a query against this arm and watch.
11. Full provenance for citation: episodes, artifacts, digests.

## Answerability map

| Question | Status |
| --- | --- |
| 3 (where evidence dies), 5 (per-case churn), 7, 8, 11 | Answerable now — served by the chunk grid, stage flow, slope graph, invalidation plan, meters, provenance |
| 1, 4, 6 | Small backend asks: verdict/score surfaced per case with worst-first ordering; text previews for case inputs and chunk content within the existing sensitivity/size policy; human descriptions authored in manifests/cases and passed through the API |
| 2, 9 | Needs judge trajectory stages (what the judge saw, its verdict text) and later a calibration projector |
| 10 | Needs a read-only playground endpoint (execute, never record) |

## Data graphics principles (Tufte, applied)

Graphics exist to answer a question faster than the table beneath them, and
each is titled with its question:

- **Paired slope graph** — "which cases moved, which regressed?" One line
  per question from baseline score to challenger score; green rises, red
  falls, gray flat; the mean is one bold line among its cases, so an average
  can never hide its outliers.
- **Chunk-presence grid** — "where did the evidence go?" Rows are chunks,
  columns are pipeline stages, filled cells are presence. A chunk that dies
  at the policy filter is visible as a row that stops.
- **Stage-flow line** — "where does the candidate set narrow?" Candidate
  counts across stages; narrowing segments drawn in red.
- **Meter bars** — "how much room do I have left?" Spent ink against the
  limit, word-sized, inline with the numbers.

Rules: direct labels, no legends, no gridlines, no decoration; graphics
word-sized or close; **missing values are never plotted as zero** — they are
excluded and named beneath the graphic. Every colored mark also carries its
word.

## Prose principles

Every section explains itself in one or two plain sentences: what this is,
why it matters ("An arm is one complete recipe for answering questions…").
The system's own descriptions (manifests, cases) should be written as prose
at authoring time and surfaced everywhere the slug appears. Terse keywords
are for machines; this tool is for a person.
