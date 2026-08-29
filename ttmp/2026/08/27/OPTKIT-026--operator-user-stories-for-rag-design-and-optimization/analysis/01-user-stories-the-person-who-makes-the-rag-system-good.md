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
RelatedFiles:
    - Path: repo://ttmp/2026/08/27/OPTKIT-028--watching-the-work-build-progress-and-the-job-view/design-doc/01-watching-a-build-progress-phases-and-the-job-view.md
      Note: Existing job-observability design whose snapshot-only history gap strengthened S38
    - Path: ws://rag-ttc/cmd/rag-ttc/cmds/experiments/answerquality/runner.go
      Note: Existing grounded-answer and judge path that motivated S42 and S49
    - Path: ws://rag-ttc/cmd/rag-ttc/cmds/indexes/build.go
      Note: Real corpus build and embedding path that motivated durable progress requirements in S38 and S50
ExternalSources: []
Summary: Fifty user stories for the person responsible for a RAG system's quality, written from needs rather than from existing surfaces — corpus knowledge, ground truth, exploration, experiments, judging, decisions, longevity, delegation, long-running work, and the end-to-end build-to-answer-benchmark workflow.
LastUpdated: 2026-08-28T21:45:00-04:00
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

Story numbers are stable identifiers, not a reading order: Theme 9 was added
after the first eight and keeps numbers S37 onward so that references from
the design documents stay valid. It belongs, in workflow order, beside
Theme 4 — it is about the long-running jobs that produce everything the
other themes look at. Theme 10 was added after reviewing the first concrete
operator goal — build and embed a large corpus, then run an LLM-backed answer
benchmark — because S21 measured retrieval support and S48 enforced job
dependencies, but neither named the end-to-end product measurement or the
single operator intention that connects the jobs.

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

## Theme 9 — Running the work

*Added after the first eight themes; belongs beside Theme 4 in workflow
order. Where Theme 4 is about designing an experiment, this theme is about
the hours during which machines are actually working — building indexes,
embedding a corpus, generating summaries, calling judges, and running
multi-turn conversations — and about my ability to see, steer, and trust
that work.*

### S37 — See everything that is running, in one place

At any moment I may have an index build going, a campaign working through
episodes, and a judge batch catching up, and all three compete for the same
API quota and the same money. Without one place that shows all of it, I
discover a stalled build only when I wonder why the campaign never started,
and I discover a runaway spend when the invoice arrives. This view is not a
dashboard for its own sake — it is how I answer "is anything wrong right
now?" in five seconds instead of five minutes.

1. Open one view listing every active piece of work, whatever kind it is.
2. For each, see what it is, what it is working on, how far along it is, how
   fast it is going, and what it has spent.
3. See the work queued behind it, and what each queued item is waiting for.
4. Drill into any one of them.
5. See recently finished work with its outcome, so "did that finish?" is
   answerable without hunting through logs.

### S38 — Watch a long build progress by phase, not by spinner

Embedding a large corpus is not one operation. It is chunking, then
representation generation, then embedding, then index construction, and each
phase has a different cost, a different failure mode, and a different
duration. A single bar across the whole thing tells me nothing I can act on:
I cannot tell whether the summarizer is slow, whether embedding is being
rate-limited, or whether I am about to exhaust the budget in a phase that
has not started yet. The phases are the granularity at which I can
intervene, so they must be the granularity I see.

1. See the phases in order, with the current one marked.
2. For each phase: units done of total, current rate, elapsed time,
   projected remaining, and cost so far.
3. See projections for phases not yet started, derived from what the
   completed phases actually measured rather than from a guess.
4. Be warned when a projection crosses a ceiling before the phase reaches it.
5. Drill into the active phase to watch individual units.
6. See a time-series graph of cumulative completion, trailing throughput,
   retries and failures, and spend, with phase boundaries and changes to
   concurrency marked on the same history.
7. Refresh, reconnect, or return after a process restart without losing that
   history; a stale last sample must look stale rather than live.

### S39 — Start, pause, stop, and resume without losing work

I need to be able to stop things. A build started with the wrong summarizer,
a campaign that is obviously going nowhere, a judge batch running against a
rubric I have just realised is broken — stopping must take effect
immediately, must not corrupt what has already completed, and must leave me
able to resume rather than restart. Anything less makes me reluctant to
start work at all, and that hesitancy is the expensive failure, not the
wasted run.

1. Start work from a description of what it will do and what it will cost.
2. Pause: in-flight units finish, nothing new is claimed.
3. Stop: the same, and the work is recorded as stopped by me rather than
   failed on its own.
4. Resume later from where it stopped, with the completed portion intact.
5. See at any moment exactly what is durable and what would be lost.

### S40 — See failures grouped by cause, while the work is still running

Long runs fail in clusters, and the cluster is the diagnosis. Four hundred
rate-limit errors mean slow down. Four hundred "document not found" mean my
inputs are wrong. Four hundred timeouts against one model mean the provider
is degraded and I should stop. A flat list of four hundred failures hides
all three readings, and I need them while the run is still going, because
two of them mean I should intervene now.

1. See failures grouped by cause, with counts, accumulating live.
2. Open a group and read the actual error from one of its members.
3. See whether a group is still growing or has stopped.
4. Retry an entire group once its cause is fixed, without re-running the
   successes.
5. Mark a group as expected so it stops demanding attention.

### S41 — Inspect one unit of work while it is in flight

When something looks wrong in aggregate, the fastest diagnosis is to look at
one unit in full: this chunk's summarization prompt and exactly what came
back, this judge call's rubric and verdict, this episode's retrieval, this
conversation's turns. Aggregates tell me that something is wrong; a single
unit tells me what it is. Without this, a misbehaving run is a statistic I
can only guess about.

1. Pick any unit from a running job.
2. See its complete inputs, including the exact prompt wherever a model is
   involved.
3. See its output, or its error, in full and unabridged.
4. See what it cost and how long it took.
5. Get from there to the configuration that produced it.

### S42 — Watch answers and judgments arrive, and stop early if the rubric is broken

An answer-and-judge batch has a property nothing else in this system has: its
output is readable prose I can evaluate immediately. If the first ten answers
are malformed, or the first ten verdicts show the judge misunderstanding the
rubric or scoring something I never asked about, then the remaining eleven
hundred calls are wasted money producing a misleading result. Reading the
actual answer and its judgment as they arrive is the cheapest quality control
available anywhere in this program, and it only works if I can see them before
the batch finishes.

1. See generated answers and verdicts streaming in as they are produced,
   newest first, while seeing how many units are waiting for generation and
   how many are waiting for judgment.
2. Open one unit and read the question, the bot's answer, its citations, the
   retrieved evidence, the judge's rationale, and the exact rubric it used.
3. Keep retrieval failure, answer-generation failure, judge failure, and a
   valid low score visibly separate in both counts and details.
4. See the score distribution forming, so a judge stuck at one value becomes
   obvious within a dozen calls.
5. Stop the batch on the spot without discarding completed answers or
   judgments.
6. Fix the answer configuration or rubric and restart, with completed work
   either reused or discarded — my choice, stated explicitly rather than
   assumed.

### S43 — Evaluate a whole conversation, not a single query

What customers actually use is a multi-turn assistant that retrieves several
times, carries context between turns, and can be right on turn one and wrong
on turn three because of what it retrieved on turn two. Scoring isolated
queries measures a component; scoring conversations measures the product.
Multi-turn evaluation needs to be a first-class kind of run with its own
results, not a script somebody wrote once and nobody can reproduce.

1. Define a scenario: an opening question, the follow-ups, and what a good
   conversation must achieve by the end.
2. Run it against a configuration, with each turn doing its own retrieval
   and tool calls.
3. See the conversation as the unit of result, with a conversation-level
   verdict.
4. See each turn underneath it, with that turn's retrieval and its own
   verdict.
5. Compare scenarios across configurations the same way single questions
   compare.

### S44 — Find the turn where a conversation went wrong

Multi-turn failures propagate. When the final answer is wrong, the cause is
usually two turns earlier: a bad retrieval that entered the context and was
never corrected. Without per-turn visibility all I know is that the
conversation failed, which tells me nothing about what to change — and
"the assistant was wrong" is not a finding I can act on.

1. Open a failed conversation.
2. See its turns in sequence, each with its verdict and what it retrieved.
3. See where the verdict first degraded, marked.
4. Open that turn's retrieval with the same explanation surfaces a single
   query gets.
5. Carry the finding out as a new question, a flag, or a proposal.

### S45 — Re-run one conversation after a change

When I fix something a single conversation exposed, I want to test the fix
against that conversation in seconds — not by re-running a suite of two
hundred and waiting an hour. The tight loop is what makes multi-turn
debugging feasible at all; without it I will avoid multi-turn work and go
back to measuring components.

1. From a failed conversation, re-run it against a different configuration.
2. See the new run beside the old one, turn by turn.
3. See which turns changed and which did not.
4. Decide whether the fix holds, and only then queue the full suite.

### S46 — Leave it running overnight and know what happened

The big runs take hours and I am not going to watch them. What I need in the
morning is not a log file but an answer: did it finish, what did it cost,
what failed and why, and is there anything that needs a decision from me
before the next step can start. Without that, every overnight run costs me
an hour of reconstruction before I can do any actual work.

1. Start the work and leave.
2. Come back to a summary: what completed, what failed grouped by cause,
   total cost against the ceiling, and duration.
3. See what stopped and needs a decision, separated from what merely
   finished.
4. Get from any line of that summary to the underlying detail.
5. See the exact corpus and bundle digests, question-suite version, answer
   model and prompt, and judge model, rubric, and epoch that produced the
   outcome.
6. Have the summary persist, so it is still there next week when someone
   asks what that run did.

### S47 — Steer throughput without editing configuration files

Concurrency and batch size are the two levers that matter during a long run.
Too aggressive and I get rate-limited, or throttled by a provider I need to
stay on good terms with. Too timid and a build takes all night for no
reason. These are runtime decisions made in response to observed behaviour,
and having to stop the run and edit a file to make them means I will simply
not make them.

1. See the current concurrency, batch size, and observed rate.
2. Change them while the work is running.
3. See the effect on throughput and error rate within a minute.
4. Have the change recorded, so a run's throughput history is part of its
   record rather than something I remember doing.

### S48 — Know what a job depends on, and be refused when the substrate is not ready

Jobs form a chain: ingest, then build, then answer benchmark, then judge.
Running a benchmark against a half-built or merely similar index produces
results that look entirely real and are not, and that is a class of mistake I
will never catch by reading numbers afterwards. The system knows the
dependency and the exact content-addressed output that satisfies it. It should
enforce both rather than let me start something that cannot be valid.

1. See, for any queued job, what it depends on and whether that dependency
   is satisfied.
2. Be refused, with the reason stated, when starting something whose inputs
   are incomplete, failed, stale, or have a different identity than the job
   declared.
3. Queue work to start automatically when the exact output it waits for is
   complete, and pass that output's immutable identity forward without my
   copying a path.
4. See a chain as one thing when it is one intention, with the phase it has
   reached and each child job's outcome.
5. Never let a downstream result omit the identities of the upstream corpus,
   bundle, question suite, answer configuration, and judge epoch.

---

## Theme 10 — Measuring the bot and chaining the work

*Added after testing the stories against the first concrete operator workflow:
build and embed a large corpus, then benchmark the answers users would
actually receive. Theme 5 deliberately separates retrieval support from judge
trust, and Theme 9 models long-running work. This theme supplies the missing
product-level measurement and the one-intention workflow connecting those
parts.*

### S49 — Run an end-to-end grounded-answer benchmark

Retrieving plausible evidence is necessary, but the customer receives an
answer, not a ranked chunk list. The answer model can omit a required fact,
misread good evidence, invent a claim, cite the wrong passage, or refuse when
it should answer. A retrieval-support score cannot see those failures. I want
the actual bot path — retrieval, prompt construction, answer generation,
citations, and judgment — measured without losing the component measurements
that explain why it behaved that way.

1. Choose a ready index bundle, a trusted question suite, the answer model and
   prompt configuration, and a versioned judge rubric.
2. Run every question through the same retrieval and answer path the product
   uses, preserving the generated answer, evidence, citations, prompts, model
   identities, latency, token usage, and cost.
3. Measure retrieval support, answer correctness and completeness, citation
   faithfulness, and refusal behaviour as separate constructs; never average
   them into one opaque quality number.
4. Keep retrieval, generation, contract-validation, and judge failures visibly
   missing rather than converting them to zero scores.
5. Read aggregate distributions and compare configurations, then drill from
   any value to the question, answer, evidence, and judge rationale behind it.

### S50 — Start one durable build-to-benchmark workflow

Building the substrate and benchmarking the resulting bot are technically
separate jobs but one intention for me: "tell me how this corpus recipe
performs." Manually copying bundle paths, remembering which build finished,
and launching the benchmark separately creates delay and makes it possible to
measure the wrong substrate. I want to describe and approve the whole chain
once, leave it working, and return to a result whose provenance is complete.

1. Select the corpus snapshot, build recipe, question suite, answer
   configuration, judge rubric, and per-resource ceilings; see phase-by-phase
   call, token, cost, and duration estimates before anything runs.
2. Approve once and start build → grounded-answer benchmark → judge → results
   as one intention with separately inspectable child jobs.
3. Automatically pass the completed bundle's immutable identity into the
   benchmark, and refuse or stop the chain when any required output is absent,
   partial, or mismatched.
4. Watch phases, progress history, throughput, spend, and grouped failures for
   the whole chain; stop at a safe boundary without losing completed work.
5. Finish at a persistent overnight summary and results view labeled with the
   corpus, bundle, question-suite, answer-model, prompt, and judge identities.

---

## What this set implies

Nine observations fall out of the stories once they are read together. They
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
- **Work itself is a first-class object the program does not have.** Theme 9
  is about builds, judge batches, campaigns, sweeps and conversation suites.
  Every one of them is long-running, budgeted, failure-prone, phased, and
  resumable, and today each is modelled separately or not at all. One job
  model — phases, units, failure causes, control, cost — would serve all of
  them; a progress surface built per kind is the expensive path and the one
  that produces five inconsistent answers to "is anything wrong?".
- **Multi-turn evaluation is a different unit, not a bigger query.** S43–S45
  need a result whose subject is a conversation with turns beneath it, each
  turn carrying its own retrieval and its own verdict. *(Refined while
  designing Theme 9: this needs a conversation executor, scenario cases, and
  projections — but no change to the measurement chain itself. Trajectory
  events already form a span tree and observation subjects are already
  generic, so a turn is a span and a turn-level score is an ordinary
  observation. See design-doc 03 §1.)*
- **Retrieval quality is not product quality.** S21 intentionally asks whether
  the evidence could support an answer; S49 asks whether the bot actually
  produced a correct, complete, grounded answer from it. Both measurements
  are needed, kept separate, because either stage can fail while the other is
  healthy.
- **A dependency chain can still be one user intention.** S48 makes child-job
  dependencies honest; S50 lets a person describe and approve the chain as a
  whole. The UI may present one workflow, but the backend should retain
  separately durable build, answer, and judge jobs so each can fail, resume,
  and be inspected truthfully.

## Consequence for ticket slicing

This document supersedes the tentative slicing proposed in OPTKIT-025 §7
(which reserved OPTKIT-026 for operations tiles and OPTKIT-027 for case
curation). That slicing was derived from the gap list rather than from user
needs; it should be re-derived from these stories once they are fleshed out,
and the operations-tile work renumbered accordingly.
