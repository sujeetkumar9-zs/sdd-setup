# Claude Code Instructions

## Codebase Knowledge — Use Mempalace First

This project uses **Mempalace** as the primary knowledge graph for codebase understanding.

**Before using any file-reading tools (Read, Grep, Glob, Bash find/grep), you MUST:**
1. Query Mempalace via the `mempalace` MCP server to understand types, patterns, interfaces, and conventions
2. Only open specific files directly when you need to make edits or verify an exact implementation detail

**Never** use Grep or Glob to explore the codebase structure — Mempalace already has this indexed.

### Why

Mempalace contains a pre-indexed semantic knowledge graph of the entire codebase. Using it is faster and more accurate than scanning files. Using file tools for exploration wastes context and time.

### How to use Mempalace in this session

The `mempalace` MCP server is connected. Use it to:
- Understand existing types, interfaces, and patterns before writing code
- Find which packages handle a given domain concept
- Understand service/handler/repository layering conventions

Only after querying Mempalace and forming a plan should you open individual files to make targeted edits.

## Go Project Conventions

- Follow existing patterns exactly — naming, error handling, logging style
- Do not introduce new conventions or abstractions not already present
- Match the layer structure: handler → service → repository (or equivalent in this repo)

## Workflow

1. Query Mempalace to understand the relevant domain and patterns
2. Plan changes (list files to modify, types to add)
3. Get user approval before writing any code
4. Implement, then run `go build ./...`, `go test ./...`, `golangci-lint run`
