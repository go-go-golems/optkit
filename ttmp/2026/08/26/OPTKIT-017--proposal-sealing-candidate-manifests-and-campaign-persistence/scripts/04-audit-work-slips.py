#!/usr/bin/env python3
from pathlib import Path
import re
import sys

slip_dir = Path(__file__).resolve().parent.parent / "various" / "work-slips"
expected = ["00-plan.log"]
for phase in range(1, 8):
    expected.extend([
        f"{phase * 2 - 1:02d}-p{phase}-start.log",
        f"{phase * 2:02d}-p{phase}-done.log",
    ])
actual = sorted(path.name for path in slip_dir.glob("*.log"))
if actual != expected:
    print(f"unexpected slip logs: expected={expected} actual={actual}", file=sys.stderr)
    raise SystemExit(1)
for name in expected:
    if not re.search(r"(?m)^printed: (?:true|yes)$", (slip_dir / name).read_text()):
        print(f"{name}: missing successful printed receipt", file=sys.stderr)
        raise SystemExit(1)
print("plan_slips: 1")
print("phase_start_slips: 7")
print("phase_done_slips: 7")
print(f"successful_print_receipts: {len(expected)}")
print("result: PASS - complete OPTKIT-017 brutalist slip sequence")
