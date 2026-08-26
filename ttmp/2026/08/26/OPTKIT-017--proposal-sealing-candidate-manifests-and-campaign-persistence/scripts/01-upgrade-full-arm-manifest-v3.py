#!/usr/bin/env python3
"""One-time strict v2 full-arm -> v3 baseline/candidate migration."""
from __future__ import annotations

import argparse
from pathlib import Path
import yaml


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--source", type=Path, required=True)
    parser.add_argument("--output", type=Path, required=True)
    args = parser.parse_args()

    source = yaml.safe_load(args.source.read_text())
    if source.get("schema") != "rag-ttc.experiment-manifest/v2":
        raise SystemExit(f"expected strict v2 source, got {source.get('schema')!r}")
    arms = source.pop("arms", None)
    if not isinstance(arms, list) or len(arms) != 2:
        raise SystemExit("migration requires exactly two full arms")
    baseline, challenger = arms
    baseline_pipeline = baseline["pipeline"]
    challenger_pipeline = challenger["pipeline"]

    expected = yaml.safe_load(yaml.safe_dump(baseline_pipeline))
    expected["retrieval"]["final_result_limit"] = challenger_pipeline["retrieval"]["final_result_limit"]
    if expected != challenger_pipeline:
        raise SystemExit("challenger differs by more than retrieval.final_result_limit")

    source["schema"] = "rag-ttc.experiment-manifest/v3"
    source["baseline"] = baseline
    source["candidates"] = [{
        "id": challenger["id"],
        **({"description": challenger["description"]} if challenger.get("description") else {}),
        "parent": baseline["id"],
        "mutations": [{
            "variable": "retrieval.final_result_limit",
            "value": challenger_pipeline["retrieval"]["final_result_limit"],
        }],
        "intent": {
            "proposer": {"kind": "human", "identity": "actor:rag-ttc-fixture"},
            "strategy": "manifest-coordinate/v1",
            "hypothesis": "Returning more fused evidence results will improve target coverage for multi-source questions.",
            "expected_improvement": {
                "metric": "retrieval.target-coverage",
                "groups": ["multi-source"],
            },
            "risks": ["additional returned evidence may include a lower-ranked distractor"],
            "motivation": {"case_ids": ["q-comparison"]},
        },
    }]

    # Keep reviewed top-level order and place authoring blocks before cases.
    ordered = {
        key: source[key]
        for key in ("schema", "name", "description", "preparation", "repeats")
        if key in source
    }
    ordered["baseline"] = source["baseline"]
    ordered["candidates"] = source["candidates"]
    ordered["cases"] = source["cases"]
    args.output.write_text(yaml.safe_dump(ordered, sort_keys=False, width=100))


if __name__ == "__main__":
    main()
