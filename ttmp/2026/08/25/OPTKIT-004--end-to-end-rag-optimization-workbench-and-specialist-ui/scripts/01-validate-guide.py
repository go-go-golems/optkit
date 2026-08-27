#!/usr/bin/env python3
"""Validate structural and repository-reference invariants in the intern guide."""

from __future__ import annotations

import re
import sys
from pathlib import Path

TICKET = Path(__file__).resolve().parents[1]
GUIDE = TICKET / "design-doc" / "01-intern-guide-to-the-rag-optimization-workbench-and-specialist-ui.md"
WORKSPACE = TICKET.parents[5]


def fail(message: str) -> None:
    raise SystemExit(f"GUIDE_VALIDATION=FAIL: {message}")


def main() -> None:
    text = GUIDE.read_text(encoding="utf-8")
    lines = text.splitlines()

    if len(lines) < 2000:
        fail(f"guide too short: {len(lines)} lines")
    if len(text.split()) < 9000:
        fail(f"guide too short: {len(text.split())} words")
    if text.count("```") % 2:
        fail("unbalanced fenced code blocks")

    headings = re.findall(r"^## (\d+)\. ", text, flags=re.MULTILINE)
    expected = [str(value) for value in range(1, 33)]
    if headings != expected:
        fail(f"numbered sections are not exactly 1..32: {headings}")

    decisions = len(re.findall(r"^### Decision:", text, flags=re.MULTILINE))
    if decisions < 8:
        fail(f"expected at least 8 decision records, found {decisions}")

    required_terms = [
        "pipeline microscope",
        "rank waterfall",
        "artifact dependency DAG",
        "complete-block",
        "Judgekit",
        "Pareto",
        "read-only",
        "source-policy",
        "pseudocode",
    ]
    lowered = text.lower()
    missing_terms = [term for term in required_terms if term.lower() not in lowered]
    if missing_terms:
        fail(f"missing required concepts: {missing_terms}")

    references = re.findall(
        r"`((?:optkit|rag-ttc|ragkit|judgekit)/[^`:\n]+\.go):\d+(?:-\d+)?`",
        text,
    )
    if len(references) < 30:
        fail(f"expected at least 30 line-anchored Go references, found {len(references)}")
    missing_files = sorted({path for path in references if not (WORKSPACE / path).is_file()})
    if missing_files:
        fail(f"missing referenced files: {missing_files}")

    print(f"GUIDE_LINES={len(lines)}")
    print(f"GUIDE_WORDS={len(text.split())}")
    print(f"GUIDE_DECISIONS={decisions}")
    print(f"GUIDE_FILE_REFERENCES={len(references)}")
    print("GUIDE_VALIDATION=PASS")


if __name__ == "__main__":
    main()
