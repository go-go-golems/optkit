#!/usr/bin/env python3
"""Per-query win/loss between two legs of a rag-ttc.index-evaluation.v1 artifact.

This is the throwaway script from OPTKIT-027, written down. It exists so that
the Go projection built in OPTKIT-029 has a target to match: the handoff's
headline number ("18 better / 11 worse / 115 unchanged across 144 queries" on
rrf) must come out of both this and the projection, or one of them is wrong.

    ./01-per-query-win-loss.py ARTIFACT [--metric ndcg_at@10] [--strategy rrf]

The comparison is per query, not per mean: a mean that moves 0.002 tells you
nothing about whether anything actually changed, and the whole point of the
OPTKIT-029 results tile is that the distribution is the finding.
"""

from __future__ import annotations

import argparse
import json
import sys
from collections import Counter


def metric_value(entry: dict, metric: str) -> float | None:
    """Read one metric off a per-query record.

    Metrics are either scalar (`mrr`) or keyed by cutoff (`ndcg_at@10`). A
    missing cutoff is None, never 0.0 — "the query was not evaluated at k=10"
    and "it scored zero at k=10" are different facts, and collapsing them is
    exactly the mistake this whole surface exists to prevent.
    """
    if "@" in metric:
        name, cutoff = metric.split("@", 1)
        table = entry.get(name)
        if not isinstance(table, dict):
            return None
        value = table.get(cutoff)
        return float(value) if value is not None else None
    value = entry.get(metric)
    return float(value) if value is not None else None


def per_query(strategy: dict, metric: str) -> dict[str, float | None]:
    return {
        entry["query_id"]: metric_value(entry, metric)
        for entry in strategy["report"].get("per_query", [])
    }


def leg(artifact: dict, bundle_index: int, strategy_name: str) -> dict:
    bundle = artifact["bundles"][bundle_index]
    for strategy in bundle["strategies"]:
        if strategy["name"] == strategy_name:
            return strategy
    raise SystemExit(f"bundle {bundle_index} has no strategy {strategy_name!r}")


def main() -> int:
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("artifact")
    parser.add_argument("--metric", default="ndcg_at@10")
    parser.add_argument("--strategy", default="rrf")
    parser.add_argument("--baseline", type=int, default=0, help="bundle index")
    parser.add_argument("--challenger", type=int, default=1, help="bundle index")
    parser.add_argument("--epsilon", type=float, default=1e-9)
    parser.add_argument("--expect", default="", help="better/worse/unchanged, e.g. 18/11/115")
    args = parser.parse_args()

    with open(args.artifact) as handle:
        artifact = json.load(handle)

    base = leg(artifact, args.baseline, args.strategy)
    challenger = leg(artifact, args.challenger, args.strategy)
    base_by_query = per_query(base, args.metric)
    challenger_by_query = per_query(challenger, args.metric)

    counts: Counter[str] = Counter()
    moved: list[tuple[str, float, float]] = []
    for query_id in sorted(set(base_by_query) | set(challenger_by_query)):
        before = base_by_query.get(query_id)
        after = challenger_by_query.get(query_id)
        if before is None or after is None:
            counts["ungraded"] += 1
            continue
        delta = after - before
        if delta > args.epsilon:
            counts["better"] += 1
            moved.append((query_id, before, after))
        elif delta < -args.epsilon:
            counts["worse"] += 1
            moved.append((query_id, before, after))
        else:
            counts["unchanged"] += 1

    baseline_id = artifact["bundles"][args.baseline]["manifest"]["bundle_id"]
    challenger_id = artifact["bundles"][args.challenger]["manifest"]["bundle_id"]
    evaluated = counts["better"] + counts["worse"] + counts["unchanged"]
    print(f"strategy   {args.strategy}   metric {args.metric}")
    print(f"baseline   {baseline_id}")
    print(f"challenger {challenger_id}")
    print(f"evaluated  {evaluated}  (ungraded {counts['ungraded']})")
    print(f"better {counts['better']}  worse {counts['worse']}  unchanged {counts['unchanged']}")
    print()
    for query_id, before, after in sorted(moved, key=lambda row: row[2] - row[1]):
        print(f"  {query_id:<24} {before:.4f} -> {after:.4f}  {after - before:+.4f}")

    if args.expect:
        better, worse, unchanged = (int(part) for part in args.expect.split("/"))
        actual = (counts["better"], counts["worse"], counts["unchanged"])
        if actual != (better, worse, unchanged):
            print(f"\nFAIL expected {better}/{worse}/{unchanged}, got {'/'.join(map(str, actual))}")
            return 1
        print(f"\nOK matches {args.expect}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
