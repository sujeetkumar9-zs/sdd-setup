# Spec-Driven Development (SDD) Skill

## Overview

You are operating in SDD mode for a Go project. Your job is to take a Jira or Confluence spec and deliver a merged PR — end to end.

## Workflow

### 1. Read the spec
- Fetch the Jira ticket or Confluence page using the Atlassian MCP
- Extract: acceptance criteria, business rules, API contracts, edge cases

### 2. Query the codebase
- Use Mempalace MCP to understand existing patterns, types, and conventions
- Identify which packages, interfaces, and files will be affected

### 3. Generate a plan
- Outline the files to create or modify
- List the interfaces/types to add
- Describe the test strategy
- **Present the plan to the user and wait for approval before writing any code**

### 4. Implement
- Follow existing code patterns exactly — naming, error handling, logging
- Write the implementation, then write tests alongside it
- Target ≥ 80% test coverage
- Run `go build ./...` mentally to check for compile errors

### 5. Quality check
- All tests must pass: `go test ./...`
- No lint errors: `golangci-lint run`
- No vet issues: `go vet ./...`
- Correctly formatted: `gofmt -l .`

### 6. PR
- Write a clear PR description referencing the Jira ticket
- Summarise what changed and why

## Rules

- Never skip the approval step after planning
- Never add code outside the scope of the spec
- Always match existing patterns — do not introduce new conventions
- If the spec is ambiguous, ask before assuming
