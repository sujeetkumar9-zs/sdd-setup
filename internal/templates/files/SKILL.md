# SDD Skill — Spec-Driven Development

You are in SDD mode. Your job: take a Jira or Confluence spec and deliver a merged PR.

## Rules

- Use the `mempalace` MCP to understand the codebase — never use Grep, Glob, or Read to explore
- Only read a file when you are about to edit it
- Never write code before presenting a plan and getting explicit approval
- Never add scope beyond the spec
- Match existing patterns exactly — naming, error handling, logging
- If the spec is ambiguous, ask before assuming

## Workflow

1. Fetch spec via Atlassian MCP
2. Search mempalace to understand affected types, patterns, and conventions
3. Plan — list files to modify and types to add, present to user, wait for approval
4. Implement following existing conventions
5. Write tests (≥ 80% coverage)
6. Run `sdd mine` to update the knowledge graph
7. Run quality gates: `go test ./...`, `go vet ./...`, `golangci-lint run`, `gofmt -l .`
8. Create PR referencing the spec
