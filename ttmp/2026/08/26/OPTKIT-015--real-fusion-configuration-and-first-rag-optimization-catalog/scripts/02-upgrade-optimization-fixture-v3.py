#!/usr/bin/env python3
"""Upgrade the frozen optimization fixture to PipelineConfig-derived v3 IDs."""

from __future__ import annotations

import hashlib
import json
from pathlib import Path

WORKSPACE = Path("/home/manuel/workspaces/2026-08-24/use-optkit")
OLD = WORKSPACE / "rag-ttc/pkg/ttc/optimization/testdata/rag-optimization-semantic-fixture-v2.json"
NEW = WORKSPACE / "rag-ttc/pkg/ttc/optimization/testdata/rag-optimization-semantic-fixture-v3.json"

IDENTITIES = {
    "corpus": "config:8cacbe6e0555b972805f6e90e71d3ef1ab2f53eae1f8871cbbdd72ed669e5415",
    "chunking": "config:5c63a7940c5523d659bd4568c647acda6aef21bfa18861a05b8f63ac5fedd460",
    "representations": "config:935a41047bac455c01413280c06e40a4c639c814ad46a780ab8fab2412f3e1e7",
    "embeddings": "config:5f3aab2f3c02eb2ff5a0d5240843bac94b271bbd06338439d604af6cc73196a7",
    "indexes": "config:c489559950e2f90c4ebdc025ad46a9dc2cfe08d4d1df0945dbeca6182a616704",
    "retrieval": "config:59ba450b5707da60d9301cc9ce937005b530a9501e728d951db9bf7eaf53b08f",
    "fusion": "config:040ca368612731f10751808236517ef58ef5d7d854136bec2305d122645c6be9",
    "reranking": "config:7cf7d081f72d5cfcf065c7c66b47a601520f54a839201840d7ec53b0bdc3c019",
    "evidence": "config:4d9a6dfcb5ff7885deefdb6031d429aa9107b664bafd8b51b230bb2583c53747",
    "context": "config:60cf498a3eb4f1e08004e22505d0dc5c52818ed2ab6e51b792a0c2125f429582",
    "answer": "config:87b0b8458db144a5c709085a56917a3ce90608831143aa3ab6244cb38d3f575e",
    "judge": "config:e51e3678bb8cb191301c8259f6fb227b20b032a15c6ab5e947d380e777cd48ab",
}

DEPENDENCIES = {
    "corpus": [],
    "chunking": ["corpus"],
    "representations": ["chunking"],
    "embeddings": ["representations"],
    "indexes": ["representations", "embeddings"],
    "retrieval": ["indexes"],
    "fusion": ["retrieval"],
    "reranking": ["fusion"],
    "evidence": ["reranking"],
    "context": ["evidence"],
    "answer": ["context"],
    "judge": ["answer"],
}

FROZEN = {"version": "fixture-v1"}
PIPELINE = {
    "corpus": FROZEN,
    "chunking": FROZEN,
    "representations": FROZEN,
    "embeddings": FROZEN,
    "indexes": FROZEN,
    "retrieval": {
        "preparation": "rag.semantic-fixture/v1",
        "route": "default",
        "final_result_limit": 2,
    },
    "fusion": {"rrf_k": 60},
    "reranking": FROZEN,
    "evidence": FROZEN,
    "context": FROZEN,
    "answer": FROZEN,
    "judge": FROZEN,
}


def main() -> None:
    data = json.loads(OLD.read_text())
    data["schema"] = "rag-ttc.optimization-semantic-fixture/v3"
    data["pipeline_config"] = PIPELINE
    for ref in data["layers"]:
        layer = ref["layer"]
        ref["identity"] = IDENTITIES[layer]
        dependencies = [IDENTITIES[name] for name in DEPENDENCIES[layer]]
        if dependencies:
            ref["depends_on"] = dependencies
        else:
            ref.pop("depends_on", None)
    data["context"]["config_identity"] = IDENTITIES["context"]
    data["answer"]["config_identity"] = IDENTITIES["answer"]
    data["judge"]["config_identity"] = IDENTITIES["judge"]
    payload = (json.dumps(data, indent=2) + "\n").encode()
    NEW.write_bytes(payload)
    OLD.unlink()
    print(f"{NEW}: sha256={hashlib.sha256(payload).hexdigest()}")


if __name__ == "__main__":
    main()
