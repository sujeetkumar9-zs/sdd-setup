Investigate and fix the bug reported in a Jira ticket.

Arguments: $ARGUMENTS
(Format: JIRA-KEY [confluence-page-id-or-url ...])

---

## Phase 1 — Context Gathering

Run all of these in parallel:

**Jira**: Fetch the issue. Extract:
- Bug description, steps to reproduce, expected vs actual behaviour
- Error messages, stack traces, or log snippets included in the ticket
- Affected environment, version, or conditions (if stated)
- Linked issues that may be related or have prior fix attempts
- Reporter and any comments that add context

**Confluence** (for each page argument): Fetch the page (or its children if empty). Extract:
- Architecture patterns, data flows, and API contracts relevant to the affected area
- Any known constraints or prior decisions that affect how the bug should be fixed

**Codebase** (mempalace):

Before querying mempalace, check for reusable context:

1. **If `.claude/sdd-contract.md` exists** — read the `## Codebase conventions` and `## Layer structure` sections. Use those values directly. Skip convention queries below.
2. **If codebase conventions were already gathered earlier in this conversation** (layer structure, error handling, logging, test/mock pattern were described in a previous skill or query result) — skip convention queries below.

**Bug-specific queries** (always run — these are specific to this ticket):
Using the service names, entity names, error messages, and code paths mentioned in the Jira ticket, search mempalace for:
- The entry point for the reported behaviour (handler, route, CLI command, event listener)
- The data flow through each layer (transport → logic → data) for the affected operation
- Types, interfaces, and services involved in the reported code path
- Any prior bug fixes or error handling patterns in the same area

**Convention queries** (skip if reusable context found above):
- Layer naming and folder structure used by this codebase
- Error handling, logging, test/mock, and dependency injection patterns

Do NOT read any file yet. Build a map of the suspect code path from mempalace before investigating further.

---

## Phase 2 — Debug Investigation

Trace the reported behaviour through the code path identified in Phase 1.

For each layer in the path:
- Identify the method(s) that handle the reported operation
- Look for: incorrect conditionals, missing error checks, wrong data transformations, off-by-one errors, race conditions, missing nil/null guards, incorrect assumptions about input shape
- Follow error propagation — check whether errors are swallowed, wrapped incorrectly, or never returned
- Check boundary conditions described in the steps to reproduce

Search mempalace as needed for each suspect location. Do NOT read files speculatively — only read a file when you have a specific reason to believe it contains the defect.

Form a root cause hypothesis. If multiple hypotheses are plausible, rank them by likelihood and note what evidence supports each.

---

## Phase 3 — RCA Report

Present the Root Cause Analysis before proposing any fix:

```
## Root Cause Analysis — <JIRA-KEY>

### Bug summary
<one sentence: what goes wrong and under what condition>

### Entry point
<the request/event/call that triggers the bug, with the method or handler name>

### Defect location
<file path>:<line range> — <method or function name>

### Root cause
<precise explanation of WHY the code produces the wrong result — not just what goes wrong, but the exact logic flaw>

### Evidence
<what in the code or the ticket's error output supports this conclusion>

### Impact
<what other callers or behaviours may be affected by this defect>

### Alternate hypotheses considered
<any other candidates that were ruled out and why — omit if none>
```

Present the RCA. Ask: "Does this root cause match what you're seeing? Any corrections or additional context?"

Revise the RCA based on user feedback if needed. **Do not proceed to Phase 4 until the user confirms the RCA is correct.**

---

## Phase 4 — Fix Plan

Based on the confirmed RCA, produce a minimal fix plan.

**The fix must be the smallest correct change.** Do not refactor surrounding code, rename variables, or improve unrelated logic. Fix the defect; nothing more.

For each change:

```
### Change <n> — <file path>
**Lines**: <line range>
**What**: <one sentence describing the change>
**Why**: <how this directly addresses the root cause>
```

After the changes, include:

```
### Regression test
**File**: <test file path>
**What**: <describe the test case that would have caught this bug — the exact scenario from the steps to reproduce>
```

If existing tests need updating to reflect corrected behaviour, list them explicitly.

Present the complete plan. Ask: "Approve this plan, suggest changes, or ask questions?"
Revise and re-present until the user gives explicit approval. **Do not write any code until approved.**

---

## Phase 5 — Implementation

For each change in the approved plan, in file order:

1. Read the file
2. Apply exactly the change described — nothing more
3. Do not modify files not referenced in the plan

After all changes, write or update the regression test described in the plan.

Run `sdd quality` and fix any failures introduced by the changes. If a pre-existing test now fails due to the fix (i.e., the test was asserting the buggy behaviour), update the test to assert the correct behaviour — do not revert the fix.

---

## Phase 6 — PR

Create a PR with:
- Title referencing the Jira key, e.g. `fix(PROJ-123): <one-line bug summary>`
- Description containing:
  - **Root cause**: one-paragraph summary from the RCA
  - **Fix**: what was changed and why
  - **Regression test**: the scenario the new/updated test covers
  - **Testing**: how to verify the fix manually (steps from the Jira ticket)

---

**Hard rules:**
- Do not proceed past Phase 3 without the user confirming the RCA
- Do not proceed past Phase 4 without explicit plan approval
- The fix is the minimum correct change — no opportunistic refactoring
- Do not read any file unless you are about to edit it or it is necessary to confirm the defect location
- If the confirmed RCA points to a defect outside the scope of what is safe to change (e.g. a shared library, a generated file, infrastructure config), surface this and ask the user how to proceed before planning
- Always include a regression test
