---
Title: 'User Stories: The Person Who Makes the RAG System Good'
Ticket: OPTKIT-026
Status: active
Topics:
    - design
    - ui
    - rag-ttc
DocType: analysis
Intent: long-term
Owners: []
RelatedFiles: []
ExternalSources: []
Summary: Thirty-six user stories for the person responsible for a RAG system's quality, written from needs rather than from existing surfaces — corpus knowledge, ground truth, exploration, experiments, judging, decisions, longevity, and delegation to an assistant.
LastUpdated: 2026-08-27T13:45:01.475999225-04:00
WhatFor: The needs-first basis from which the workbench's remaining surfaces and ticket slicing are derived.
WhenToUse: Read before designing any new surface; use as the acceptance vocabulary for OPTKIT-026 and successors.
---

# User Stories: The Person Who Makes the RAG System Good

## Who this is for

The primary user is the person accountable for whether TTC's retrieval
actually works. They did not write the retrieval code and will not read it.
They may not own the knowledge base either — content is often someone else's
job — but they are the person who has to say "this configuration is better
than that one, and here is why", and to be right about it. They have a
budget, a limited number of hours, and an AI assistant that can do tedious
work if supervised.

Two secondary users appear throughout: the **assistant**, which proposes and
prepares but must not commit; and the **inheritor** — a colleague, a new
hire, or the same person six months later — who has to understand why the
system is configured the way it is without being told the story in person.

## How to read this

Each story has two paragraphs: **what and why** (the need, and what goes
wrong without it), then **how** (the user-visible steps). The steps are
deliberately written in the user's terms, not in terms of anything that
exists today: no story here names a tile, an endpoint, or a command. That
is the point — this document is the needs-first basis, and what we build is
derived from it rather than the other way around.

These stories are a first overview. Each will be fleshed out with concrete
acceptance criteria, data requirements, and surface design later.

---

## Theme 1 — Knowing my material

### S1 — See what is actually in my corpus

Before deciding anything about chunking, representations, or retrieval, I
need a truthful picture of what I am indexing: how many documents, of what
kinds, how long, how uniform, and how much of it is navigation boilerplate
or the same marketing paragraph repeated across forty pages. Every
downstream decision depends on facts about the material, and right now those
facts live in my head as guesses. Getting this wrong is expensive in a
specific way: I can spend a week tuning fusion when the real problem is that
a third of the corpus is near-duplicate product blurbs.

1. Point the system at the raw source material.
2. See a profile: counts by type and source, a length distribution,
   near-duplicate clusters, documents with almost no substantive text.
3. Drill into any bucket and read the actual documents behind the numbers.
4. Mark what should be excluded and why, so the exclusion becomes part of
   the corpus definition rather than a filter I re-apply from memory.

### S2 — Find out whether my corpus can answer the questions people ask

Retrieval can only surface what exists. If customers ask about overwintering
a plant we never wrote about, no configuration fixes that. I need to
separate a retrieval problem from a content problem early, because they
route to different people: one is my job, the other belongs to whoever
writes content. Conflating them wastes weeks and produces optimization
pressure on a system that is already doing the best it can.

1. Bring a list of real questions.
2. For each, find out whether any document plausibly answers it — searched
   generously, by any means available, not by one configuration.
3. Get a partition: answerable, partially answerable, not covered.
4. Send the "not covered" list to the content owner as a work item.
5. Keep the rest as the optimization target, and revisit the partition as
   content is added.

### S3 — Record what I know is junk, duplicated, or missing

My knowledge of the corpus's flaws is currently a set of things I remember.
If it is not recorded, every future run re-learns it, my assistant cannot
use it, and the next person starts from zero. I want corpus-level
annotations that persist and that show up wherever the annotated thing
appears later.

1. While reading documents or results, flag an item: junk, duplicate of X,
   outdated, restricted.
2. Give a reason in a sentence.
3. Choose whether the flag excludes the item from indexing or is only a note.
4. See the flag surface automatically anywhere that document appears again,
   including inside future results.

---

## Theme 2 — Defining what "good" means

### S4 — Assemble a question set that reflects reality

The evaluation set *is* the definition of success. A set of questions I
invented at my desk will optimize the system for my imagination. I want
questions drawn from actual traffic, support tickets, and the hard cases I
already know about, with enough variety — short and long, factual and
comparative and procedural — that a change which helps one kind and hurts
another becomes visible instead of averaging out.

1. Import questions in bulk from logs, tickets, or a spreadsheet, and add
   my own by typing them.
2. See them as a working set with categories and counts.
3. Collapse near-identical questions.
4. Mark which ones matter most, by frequency or business consequence.
5. Grow the set over time without invalidating past results: additions are
   versioned rather than silent.

### S5 — Say what a right answer looks like, question by question

A score means nothing unless "correct" is written down somewhere I can
read. For each question I want to name the documents or passages that ought
to be found, and to write in words what a good retrieval must support. The
prose matters as much as the target list: it is what a judge reads, what a
colleague reads, and what I read when I have forgotten why I picked those
targets.

1. Open a question.
2. Attach the documents that ought to be retrieved, each with a grade —
   essential, or merely helpful.
3. Write one or two sentences: what must this evidence let someone say?
4. Note explicit non-answers: documents that look relevant and are not.
5. Save; the question is now scoreable, by a metric and by a judge.

### S6 — Author ground truth from evidence in front of me

I cannot write correct target lists from memory. I do not know the document
identifiers, and I do not know the corpus well enough to name what should
match. The only practical method is: ask the question, look at what comes
back, and mark the good ones. This turns ground-truth authoring from data
entry into review, which is the difference between a suite of fifteen
questions and a suite of a hundred and fifty.

1. Ask a question against any configuration — deliberately a generous one,
   to see more candidates than a tuned one would return.
2. Read the returned passages.
3. Mark each: essential, helpful, irrelevant, actively misleading.
4. Search differently to find targets retrieval missed, and add those too.
5. The marks become the question's judgments, with a record of how I made
   them.

### S7 — Change my mind about ground truth without corrupting history

I will discover that a question was badly specified, or that a document I
marked irrelevant is in fact the best answer. I need to fix it. I also need
past results to stay interpretable — either re-scored against the corrected
truth, or clearly labeled as measured under the old one. What I must never
have is a chart that silently mixes both.

1. Edit a question's targets or rubric.
2. See which past results were measured under the previous definition.
3. Choose: re-score the affected results, or keep them labeled as measured
   under the old definition.
4. The change is recorded with who changed it, when, and why.

### S8 — Keep a small set of questions I trust absolutely

I need a subset I have personally verified in detail — gold cases — to
check both the system and the judge. When a metric moves, gold cases are
where I look first, and they are the only material against which "the judge
agrees with me" can be measured at all.

1. Mark a question gold after reviewing it end to end.
2. Record my own score or verdict for it.
3. See gold cases called out separately in every result view.
4. Be told when the gold subset and the full set disagree about the
   direction of a change.

---

## Theme 3 — Building intuition

### S9 — Ask a question right now and see what comes back

This is the single most-used action in the entire workflow, and it is not a
measurement — it is how I learn what the system does. I want to type a
question, choose which configuration answers it, and see results
immediately, freely enough that I can try five phrasings in a minute
without any of it counting as an experiment or costing me a decision.

1. Type a question.
2. Pick a configuration, or accept the current one.
3. See ranked results with their text, source, and score, fast.
4. Re-ask with a variation.
5. Either keep the query — promoting it into the question set — or discard
   it. Nothing here enters the record unless I say so.

### S10 — Understand why a result ranked where it did

A ranked list without explanation teaches nothing. When something obviously
good sits at rank nine, I need to see whether keyword search or vector
search found it, what text it matched, what each channel scored it, and
where fusion placed it. That is the difference between guessing at a fix
and knowing which decision is implicated.

1. Pick a result.
2. See its path: which retrievers found it, at what rank and score, which
   text it matched on, how fusion combined those, and whether any filter
   nearly dropped it.
3. See the competitors that beat it, with the same detail.
4. Move from there directly to the setting that governs the step that hurt
   it.

### S11 — See how a document is split under different chunking strategies

Chunking has the widest blast radius and the least visibility of any
decision I make. Aggregate metrics can tell me strategy A beat strategy B;
they cannot tell me why, and they do not help me choose what to try next.
Seeing one real document cut both ways — where the boundaries fall, which
chunk keeps the heading, what got orphaned mid-list — is how I form a
hypothesis worth spending money on.

1. Pick a document I know well.
2. Choose two or more chunking settings.
3. See the document rendered side by side with boundaries marked, chunk
   sizes, and which chunks lost their heading context.
4. Ask a question and see which chunks each strategy would have surfaced.

### S12 — See what the retriever actually matched on

Once I add presummarization, the searchable text is no longer the document
text. When something matches, I need to know whether it matched raw text, a
generated summary, or a generated question — otherwise I cannot tell
whether the summarizer is helping or quietly rewriting my corpus into
something that retrieves well and cites badly.

1. In any result, reveal the text that was actually matched.
2. See it labeled by kind — raw, summary, question, entities — beside the
   source passage it was derived from.
3. Flag a derived text that misrepresents its source.
4. Have that flag count against the approach that produced it, not just sit
   as a note.

### S13 — Put two configurations side by side on one question

Before committing to a full measurement I want a cheap read on whether a
change does anything at all. One question, two configurations, differences
highlighted, answers "is this worth measuring properly?" in seconds — and
saves me from running a two-hour campaign to discover a knob was inert.

1. Ask a question.
2. Add a second configuration to answer the same question.
3. See both lists aligned, with moved, new, and dropped results marked.
4. Repeat across a handful of questions to see whether the pattern holds.
5. If it looks promising, turn this exact pair into a real experiment
   without re-specifying it.

---

## Theme 4 — Running experiments

### S14 — Try several variants and learn which is best

This is the core act of the job. I have a hypothesis — smaller chunks help
procedural questions — and I want it measured across the whole question set,
repeated enough times to be believable, without babysitting the process.
The output I want is not a number but a ranking with enough per-question
detail to argue about.

1. Start from a configuration that exists.
2. Describe what varies: a setting, a list of values, a substrate change.
3. Confirm the resulting set of variants and what running them will cost.
4. Start it, and leave.
5. Come back to a ranked comparison with per-question detail underneath.

### S15 — Know what a run will cost and how long it will take, before it runs

Some of these runs cost real money — embeddings, summarization, judging —
and hours of wall-clock time. I want an estimate that names each resource
before anything is spent, and a ceiling I set that the system will not
cross. An unpleasant surprise here does not just cost money; it makes me
reluctant to experiment, which is the actual damage.

1. After describing the experiment, see a table: how many runs, how many
   embedding tokens, how many judge calls, estimated duration, and cost in
   currency where it is known.
2. Adjust the scope — fewer questions, fewer repeats, a cheaper judge — and
   watch the estimate change.
3. Set ceilings per resource.
4. Approve the spend explicitly.
5. If a ceiling is reached mid-run, the run pauses and asks rather than
   continuing or dying.

### S16 — Watch a run and intervene

Long runs go wrong in ways that are obvious ten minutes in: everything
failing, one variant timing out, spend accelerating faster than projected.
I want to see progress, failures with their reasons, and burn rate while it
happens, and to stop it without losing the work already completed.

1. See live counts — waiting, running, done, failed — and throughput.
2. See failures with their reason as they occur, not in a summary at the end.
3. Inspect any single run in flight or just finished.
4. Pause or stop. Completed work is kept and remains valid.

### S17 — Resume or repair a run that broke

Provider outages, rate limits, and my laptop closing are all normal events.
Losing four hours of work to any of them is not acceptable. Neither is a
half-finished run that reports itself as complete — that is worse than
losing the work, because I would act on it.

1. See clearly that a run is incomplete, and why.
2. Retry only the failed portion.
3. See which failures are permanent — a malformed question, a missing
   document — and which are transient.
4. Either complete the run, or mark it partial, with results that state
   their partiality everywhere they appear.

### S18 — Change the substrate as part of an experiment

Chunking, presummarization, embedding model, and index type are the
decisions I most want to test, and they are exactly the ones that require
rebuilding an index. I do not want to hand-manage builds and remember which
index went with which result — that is how results get attributed to the
wrong cause. I want to name a substrate variation as a dimension of an
experiment and have the machinery work out what must be rebuilt and what
can be reused.

1. Describe the variation: a different chunk size, add summaries, a
   different embedding model.
2. See what that invalidates and what is reused, with the cost of the
   difference.
3. Approve.
4. The rebuild runs as part of the experiment, visible like any other work.
5. Results are labeled with the substrate they were measured on,
   permanently and without my having to remember.

### S19 — Sweep a range without hand-building each variant

Some questions are "what is the best value of X", not "is A better than B".
Authoring nine near-identical variants by hand is boring and error-prone,
and a typo in one of them produces a result I will believe.

1. Pick a setting.
2. Give a range or a list of values.
3. See the variants generated, with the total cost.
4. Run them.
5. See the result as a curve with a visible optimum, not a table of pairs I
   have to read like a spreadsheet.

### S20 — Only compare things that are comparable

The most dangerous failure mode in this whole system is a number that looks
like progress but was produced differently: a different judge model, a
different question set, a different corpus. I want the system to refuse
silently-invalid comparisons rather than render them attractively. I will
not catch this myself — that is precisely why it is dangerous.

1. When I select two results to compare, be told immediately if they differ
   in what was measured or how.
2. See exactly what differs.
3. Either compare only the subset that is genuinely comparable, or re-measure
   one side to match the other.
4. Never be shown a delta computed across a difference that invalidates it.

---

## Theme 5 — Judging quality

### S21 — Measure whether the evidence would actually support an answer

Overlap with a target list is cheap, reproducible, and blind: it cannot
distinguish "found the right page" from "found enough to answer the
question". For a customer-facing assistant, the second is the thing that
matters, and at any scale only a language model can assess it. I want that
assessment as a second measurement standing beside the first, never blended
into it.

1. Define, per question, what a good answer must contain.
2. Have every run assessed by a judge against that definition.
3. See the judge's score beside the overlap score for every question.
4. Treat them as two distinct measurements — never averaged into a single
   headline number.

### S22 — Read the judge's reasoning and see what it saw

A score of 0.6 from a model is not evidence; the reasoning is. I need to
read why it scored that way and see the exact evidence text it was given.
Without that I cannot tell a genuine weakness in retrieval from an artifact
of how the judge was prompted, and I will end up either trusting it blindly
or ignoring it entirely.

1. Open any judged result.
2. Read the verdict: per-dimension scores and the written rationale.
3. See the exact evidence that was presented to the judge.
4. See which model and which prompt produced it.
5. Flag it when it is wrong.

### S23 — Overrule the judge, and have it count

When I disagree with the judge, my opinion should become data rather than a
complaint. My labels are the only ground truth for judge quality that will
ever exist, and they should accumulate as a by-product of work I am already
doing — not in a separate labeling exercise I will schedule and never do.

1. On any judged result, record my own verdict and a sentence of reasoning.
2. See it stored against that question and that judge version.
3. See my labels feed the judge's agreement statistics automatically.
4. Have my label take precedence wherever both mine and the judge's exist.

### S24 — Know whether the judge can be trusted right now

A judge that has drifted, or that never agreed with me on this kind of
question, makes every number downstream fiction. I want a standing answer
to "how much do I trust this judge?" — not a one-time validation I did once
and then forgot while the model was silently updated underneath me.

1. See agreement between the judge and my labels, overall and by question
   type.
2. See where they disagree most, with examples I can open.
3. Be warned when agreement drops after a judge or prompt change.
4. Re-judge past results under a new judge version and compare the two
   judgments directly.

### S25 — Tell "failed" apart from "scored badly"

An API timeout is not a zero. If failures silently become low scores I will
chase phantom regressions, and once I discover that has happened I will
stop trusting every number the system has ever shown me. Missing must stay
visibly missing, at every level of aggregation.

1. See failures counted separately from scores, everywhere.
2. See the reason each one failed.
3. See how much of a result is missing before I read its average.
4. Retry the failures without re-running everything else.

---

## Theme 6 — Deciding

### S26 — See what got worse, not just what got better

An average improvement almost always hides a set of questions that
regressed, and those are frequently the ones customers notice. Before
adopting anything I want the losers listed, worst first, with enough detail
to judge whether the trade is one I am willing to make and defend.

1. Compare the candidate against what is current.
2. See per-question deltas sorted by damage.
3. Open any regression and see exactly what changed in what was retrieved.
4. Decide whether the trade is acceptable, and record that reasoning where
   it will be found later.

### S27 — Know whether a difference is real

With thirty questions and a stochastic judge, small deltas are noise. I want
to know when a result is solid enough to act on, expressed in a way I can
actually act on — not a statistic I will misread, and not a false precision
that makes a coin flip look like a finding.

1. See the spread across repeats and across questions, not only the mean.
2. See how many questions moved in each direction.
3. Be told plainly when a difference is within noise.
4. Be offered the cheapest way to raise confidence: more repeats, more
   questions, or a tighter comparison.

### S28 — Record the decision and what justified it

In three months I will not remember why fusion is set the way it is, and
neither will anyone else. The decision — what I chose, what I rejected, the
trade-off I accepted, and who made the call — is the most valuable artifact
this entire process produces, and it is the one most likely to be lost.

1. After comparing, write the decision: what I chose, what I rejected, why.
2. Attach the results and the specific cases that convinced me.
3. Note the known risks and what I would watch for.
4. Have the decision appear whenever anyone looks at that configuration
   afterwards, without having to go looking for it.

### S29 — Adopt a configuration and know what it changes

Moving a chosen configuration into the live system is where risk becomes
real. I want to see precisely what differs from what is running, what must
be rebuilt or redeployed, and what could break — before I commit, and with
a way back that I have confirmed exists rather than assumed.

1. Select the winning configuration.
2. See a difference against what is live, and what that difference
   invalidates.
3. See the rebuild and deployment consequences, with cost and duration.
4. Approve.
5. Keep the previous configuration reachable, so reverting is a decision
   rather than a recovery project.

---

## Theme 7 — Living with the system

### S30 — Re-check after the corpus changes

Content is added and edited constantly, and a configuration tuned on last
month's corpus can degrade quietly. I want the same evaluation re-run on
the new corpus and to be told what moved, without redesigning the
experiment or wondering whether I ran it the same way.

1. After a content update, re-run the standing evaluation.
2. See the comparison against the last run over the same questions.
3. See which regressions come from content changes and which from
   configuration.
4. Route the content problems to the content owner with the evidence
   attached.

### S31 — Reproduce an old result

When a number is questioned — by me, by a colleague, or during an incident
— I need to show exactly what produced it: which corpus, which chunking,
which model, which judge, which questions. Without that, results are
anecdotes, and anecdotes lose arguments to whoever speaks with more
confidence.

1. Open any past result.
2. See the complete identity of everything that produced it.
3. Verify that the record is intact and unaltered.
4. Re-run it as it was, if the inputs still exist, and see whether it
   reproduces.

### S32 — Know what the system costs to run

I am accountable for spend. I want to know what experimentation cost this
month, what a production query costs, and which of my choices —
summarization, the judge, the embedding model — are driving it. Cost is
also a legitimate tiebreaker: two configurations that score the same are
not equally good if one costs three times as much per query.

1. See spend per experiment and cumulatively over time.
2. See cost attributed to each kind of activity.
3. See the estimated per-query production cost of a configuration.
4. Use that cost as an input when choosing between configurations that
   perform similarly.

### S33 — Hand the story to someone else

Someone will inherit this: a colleague, a new hire, or me after six months
on something else. The sequence of what was tried, what was learned, and
what was decided is the real accumulated value, and it should be readable
without my narration. If it only exists in my head, the program's output is
a configuration file and nothing else.

1. Open a chronological view of the work: experiments, decisions, corpus
   changes.
2. Read each decision with the evidence that supported it.
3. Follow any claim back to the underlying results.
4. Export or share it as a document that stands on its own.

---

## Theme 8 — Working with an assistant

### S34 — Delegate the tedious parts

Drafting eighty candidate questions, proposing a sweep, triaging forty
failures into groups, writing the first version of a rubric — this is work
an assistant does well and I do slowly. I want to hand it over in my own
words, referring to whatever I happen to be looking at, and get back
something concrete enough to review rather than a suggestion I still have
to implement.

1. Describe what I want in prose, referring to things in front of me.
2. The assistant proposes concrete work: draft questions, a variant set, a
   triage summary.
3. I review it in the same place I would have authored it myself.
4. Accept, edit, or reject item by item rather than all at once.

### S35 — Approve anything that costs money or becomes permanent

I want the assistant to be able to prepare anything and commit nothing.
Spending money, adopting a configuration, or writing something durable
should require my explicit approval each time, with the consequence stated
in the request. This is what lets me give it more rope: the blast radius of
a mistake is bounded by what I approved.

1. The assistant requests an action.
2. I see what it is, what it costs, and what it changes permanently.
3. I approve it, modify it first, or refuse.
4. The record shows that the assistant proposed and that I approved — both
   facts, durably.

### S36 — See and undo what the assistant did

Trust requires audit. I want a plain record of every action the assistant
took, distinguishable at a glance from my own, and the ability to reverse
anything that is not already permanent. Without the undo, I will supervise
too tightly to get any leverage from the delegation.

1. See a chronological log of actions with who performed each.
2. Filter it to just the assistant's.
3. Undo the reversible ones.
4. See clearly which ones are permanent, and why they cannot be undone.

---

## What this set implies

Five observations fall out of the stories once they are read together. They
are noted here as consequences to test during fleshing-out, not as design
decisions.

- **An interactive query path is load-bearing.** S9–S13 are the exploration
  loop, and S6 — authoring ground truth from evidence — depends on it
  completely. A large fraction of the other stories are unreachable without
  the ability to ask an arbitrary question and look at the answer.
- **Ground-truth authoring is the hinge.** S5–S8 are where exploration turns
  into measurement. They are also the most human-hours-expensive stories in
  the set, which makes their ergonomics disproportionately important.
- **Substrate variation is not an advanced feature.** S11, S12, S18 and S30
  all want the index-time decisions — chunking, presummarization, embedding
  — treated as ordinary experimental dimensions. That is the main thing the
  user wants to study, not a later refinement.
- **The honesty stories are mostly about surfacing guarantees.** S20, S25,
  S27 and S31 describe properties the record layer already promises. The
  work is making those properties legible at the moment a person is about
  to act on a number.
- **Decision records are the actual product.** S28 and S33 describe the
  artifact with the longest useful life, and it is the one with no home
  today. Everything else produces evidence; these produce the conclusion.

## Consequence for ticket slicing

This document supersedes the tentative slicing proposed in OPTKIT-025 §7
(which reserved OPTKIT-026 for operations tiles and OPTKIT-027 for case
curation). That slicing was derived from the gap list rather than from user
needs; it should be re-derived from these stories once they are fleshed out,
and the operations-tile work renumbered accordingly.
