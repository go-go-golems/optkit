#!/usr/bin/env python3
from pathlib import Path
import re
import sys

script_dir = Path(__file__).resolve().parent
slip_dir = script_dir.parent / "various" / "work-slips"
expected = ["00-plan.log"]
for phase in range(1, 8):
    start_number = phase * 2 - 1
    done_number = phase * 2
    expected.extend([
        f"{start_number:02d}-p{phase}-start.log",
        f"{done_number:02d}-p{phase}-done.log",
    ])

actual = sorted(path.name for path in slip_dir.glob("*.log"))
if actual != expected:
    print(f"unexpected slip log set:\nexpected={expected}\nactual={actual}", file=sys.stderr)
    raise SystemExit(1)
for name in expected:
    text = (slip_dir / name).read_text()
    if not re.search(r"(?m)^printed: (?:true|yes)$", text):
        print(f"{name}: missing successful printed receipt", file=sys.stderr)
        raise SystemExit(1)
print(f"plan_slips: 1")
print(f"phase_start_slips: 7")
print(f"phase_done_slips: 7")
print(f"successful_print_receipts: {len(expected)}")
print("result: PASS - complete OPTKIT-016 brutalist slip sequence")
