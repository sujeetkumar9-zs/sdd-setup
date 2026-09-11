# SDD Skill — Address PR Comments

You are in address-pr-comments mode. Your job: take all unresolved reviewer comments on the current branch's open PR, plan the changes needed to satisfy them, get approval, implement them, and record what was addressed so future runs don't re-do the same work.

## Rules

- Use `mempalace` MCP for all codebase knowledge — never use Grep, Glob, or Read to explore
- Only read a file when you are about to edit it
- Never write code before presenting a plan and getting explicit approval
- Never add scope beyond what the reviewer comments request
- Match existing patterns exactly — naming, error handling, logging
- History file (`.claude/pr-feedback-history.md`) is append-only — never overwrite past entries
- If a comment conflicts with Jira acceptance criteria, surface the conflict before acting

## Inputs

- **PR**: inferred from current branch (fetched via GitHub MCP)
- **JIRA-KEY** (optional): acceptance criteria and scope guard
- **Confluence pages** (optional): architecture constraints
- **`.claude/pr-feedback-history.md`** (optional): previously addressed comments — skip these

## Comment Classification

Every unresolved comment is one of:
- **ACTIONABLE** — clear request, known change, proceed to plan
- **AMBIGUOUS** — unclear intent or scope → ask all questions in one batch before planning
- **OUT OF SCOPE** — beyond Jira acceptance criteria → surface to user, do not implement without explicit approval

## Workflow

### Phase 1 — Context (parallel)
- Fetch PR diff + unresolved review comments
- Fetch Jira issue (if provided)
- Fetch Confluence pages (if provided)
- Load history file; cross-reference to skip already-addressed comment IDs
- Search mempalace for patterns in files touched by the diff

### Phase 2 — Triage
- Classify all comments: ACTIONABLE / AMBIGUOUS / OUT OF SCOPE
- Batch all clarifying questions for AMBIGUOUS comments and ask at once
- Surface OUT OF SCOPE comments and ask whether to include
- Do NOT proceed until all ambiguous comments are resolved or deferred

### Phase 3 — Plan
For each ACTIONABLE comment: file path, line range, exact change, rationale.
Present summary table. Get explicit approval. No code until approved.

### Phase 4 — Implementation
- Read → Edit each file in the approved plan
- Exactly the changes described — no broader cleanup
- Run `sdd quality`; fix any failures introduced by the changes

### Phase 5 — History + Confirm
- Append session block to `.claude/pr-feedback-history.md` (addressed, deferred, out of scope)
- Summarise to the user: comments addressed, files changed, anything skipped and why
