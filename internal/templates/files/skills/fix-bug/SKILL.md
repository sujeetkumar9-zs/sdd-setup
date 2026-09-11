# SDD Skill — Fix Bug

You are in fix-bug mode. Your job: investigate the bug described in a Jira ticket, produce a Root Cause Analysis, get it confirmed, plan the minimal fix, get it approved, implement it, and create a PR.

## Rules

- Use `mempalace` MCP for all codebase knowledge — never use Grep, Glob, or Read to explore
- Only read a file when you are about to edit it, or when you have a specific reason to believe it contains the defect
- Never write code before presenting a plan and getting explicit approval
- **Never proceed past the RCA without user confirmation** — fixing the wrong root cause is worse than not fixing at all
- The fix must be the minimum correct change — no opportunistic refactoring, no unrelated cleanup
- Always include a regression test
- Match existing patterns exactly — naming, error handling, logging
- If the defect is in a shared library, generated file, or infrastructure config, surface this before planning

## Inputs

- **JIRA-KEY** (required): bug description, steps to reproduce, expected vs actual, error messages
- **Confluence pages** (optional): architecture docs, data flows, known constraints

## Workflow

### Phase 1 — Context (parallel)
- Fetch Jira ticket: bug description, repro steps, expected vs actual, stack traces, linked issues
- Fetch Confluence pages (if provided): architecture, data flows, API contracts for the affected area
- Search mempalace: entry point, data flow through all layers, types/interfaces/services in the suspect path
- Build a map of the code path before reading any file

### Phase 2 — Debug Investigation
- Trace the reported behaviour through each layer of the code path
- Look for: incorrect conditionals, swallowed errors, wrong data transformations, missing nil guards, race conditions, boundary condition failures
- Search mempalace for each suspect location — read a file only when you have a specific reason to believe it contains the defect
- Form a ranked list of root cause hypotheses; identify the most likely

### Phase 3 — RCA Report
Present a structured RCA:
- Bug summary (one sentence)
- Entry point (method/handler that triggers it)
- Defect location (file:line + method name)
- Root cause (the exact logic flaw, not just the symptom)
- Evidence (what in the code or ticket supports this)
- Impact (other callers or behaviours affected)
- Alternate hypotheses considered (and why ruled out)

**Ask the user to confirm the RCA. Do not proceed until confirmed.**

### Phase 4 — Fix Plan
- One entry per change: file path, line range, what changes, why it addresses the root cause
- Include a regression test: the exact scenario from the repro steps
- List any existing tests that need updating (test was asserting buggy behaviour)
- Get explicit approval. No code until approved.

### Phase 5 — Implementation
- Read → Edit each file in the approved plan
- Exactly the changes described — nothing more
- Write or update the regression test
- Run `sdd quality`; fix any failures; update tests asserting buggy behaviour rather than reverting the fix

### Phase 6 — PR
- Title: `fix(JIRA-KEY): <one-line bug summary>`
- Description: root cause summary, what changed and why, regression test scenario, manual verification steps from the Jira ticket
