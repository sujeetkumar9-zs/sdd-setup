Address reviewer comments on the open PR for the current branch.

Arguments: $ARGUMENTS
(Format: [JIRA-KEY] [confluence-page-id-or-url ...])
All arguments are optional. The PR is inferred from the current branch.

---

## Phase 1 — Context Gathering

Run all of these in parallel:

**PR** (GitHub MCP): Fetch the open PR on the current branch. Extract:
- PR title, description, diff (all changed files and lines)
- All review comments with: comment ID, author, file path, line number, exact text, and resolution status
- Filter to UNRESOLVED comments only — skip any already marked resolved on GitHub

**Jira** (if JIRA-KEY provided): Fetch the issue. Extract:
- Acceptance criteria and any constraints
- Linked issues that may affect scope

**Confluence** (for each page argument): Fetch the page (or its children if empty). Extract:
- Architecture patterns, API contracts, field definitions
- Any constraints relevant to the changed code

**History**: Load `.claude/pr-feedback-history.md` if it exists. This file records comment IDs that were addressed in previous runs of this skill on this PR. Cross-reference against the PR's open comments — skip any comment whose ID appears in the history as addressed.

**Codebase** (mempalace):

Before querying mempalace, check for reusable context:

1. **If `.claude/sdd-contract.md` exists** — read the `## Codebase conventions` and `## Layer structure` sections. Use those values directly. Skip convention queries below.
2. **If codebase conventions were already gathered earlier in this conversation** (layer structure, error handling, logging, test/mock pattern were described in a previous skill or query result) — skip convention queries below.

**PR-specific queries** (always run — these are specific to this PR):
For each file touched by the PR diff, search mempalace for:
- The layer it belongs to and the patterns used there
- Patterns directly relevant to the reviewer's concerns (e.g. if a comment is about error handling, search for the error handling pattern in that layer)

**Convention queries** (skip if reusable context found above):
- General error handling, naming, test/mock conventions used across the codebase

---

## Phase 2 — Triage

Classify every unresolved, non-history comment into exactly one bucket:

**ACTIONABLE** — the request is clear; you know what to change and where
**AMBIGUOUS** — the intent is unclear, the scope is undefined, or multiple valid interpretations exist
**OUT OF SCOPE** — the comment requests changes beyond the Jira acceptance criteria (flag, do not implement without approval)

For AMBIGUOUS comments: collect all clarifying questions and present them as a single batch to the user before doing anything else. Wait for answers. Do NOT proceed to Phase 3 until all ambiguous comments are resolved or explicitly deferred by the user.

For OUT OF SCOPE comments: list them with a note explaining why they are out of scope. Ask the user whether to include them anyway before adding to the plan.

---

## Phase 3 — Plan

For every ACTIONABLE comment (plus any AMBIGUOUS ones resolved in Phase 2), produce a plan entry:

```
### Comment #<id> — <author> on <file>:<line>
> "<reviewer's exact comment text>"

**Change**: <one-sentence description of what will change>
**File**: <exact path>
**Lines**: <line range or "new code">
**Rationale**: <why this satisfies the comment, tied to the acceptance criterion if applicable>
```

Group by file for readability. At the end of the plan, show a summary table:

| # | File | Author | Status |
|---|------|--------|--------|
| <id> | <file:line> | <author> | ACTIONABLE / DEFERRED / OUT OF SCOPE |

Present the complete plan. Ask: "Approve this plan, suggest changes, or ask questions?"
Revise and re-present until the user gives explicit approval. **Do not write any code until approved.**

---

## Phase 4 — Implementation

For each comment in the approved plan, in file order:

1. Read the file
2. Make exactly the change described in the plan — no broader refactors, no extra cleanup
3. Do not modify files not referenced in the plan

After all changes are made, run `sdd quality` and fix any failures introduced by the changes before proceeding.

---

## Phase 5 — History Update and PR

Update `.claude/pr-feedback-history.md`:

```markdown
# PR Feedback History

## PR: <PR-number> — <PR-title>
<!-- append; do not overwrite previous entries -->

### Session: <today's date>
Addressed comments:
- #<comment-id> (<file>:<line>) — <one-line summary of change made>
- #<comment-id> ...

Deferred:
- #<comment-id> — <reason deferred>

Out of scope (not addressed):
- #<comment-id> — <reason>
```

If `.claude/pr-feedback-history.md` already exists, **append** the new session block — do not overwrite previous entries.

Finally, confirm to the user:
- How many comments were addressed
- Which files were changed
- Which comments were deferred or skipped and why
- Remind them to push the branch and re-request review

---

**Hard rules:**
- Do not write any code before Phase 3 approval
- Do not modify files not referenced in the approved plan
- Do not address a comment already recorded as addressed in `.claude/pr-feedback-history.md`
- Do not add scope beyond what the comments request
- If a comment conflicts with the Jira acceptance criteria, flag the conflict before implementing
- History file is append-only — never delete or overwrite past entries
