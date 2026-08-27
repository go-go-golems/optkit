# Propose Workspace: Draft Document and RRF Vertical Slice

This is the document workspace for ticket OPTKIT-023.

## Structure

- **design/**: Design documents and architecture notes
- **reference/**: Reference documentation and API contracts
- **playbooks/**: Operational playbooks and procedures
- **scripts/**: Utility scripts and automation
- **sources/**: External sources and imported documents
- **various/**: Scratch or meeting notes, working notes
- **archive/**: Optional space for deprecated or reference-only artifacts

## Getting Started

Use docmgr commands to manage this workspace:

- Add documents: `docmgr doc add --ticket OPTKIT-023 --doc-type design-doc --title "My Design"`
- Import sources: `docmgr import file --ticket OPTKIT-023 --file /path/to/doc.md`
- Update metadata: `docmgr meta update --ticket OPTKIT-023 --field Status --value review`
