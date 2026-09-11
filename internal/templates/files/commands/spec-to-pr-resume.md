Resume an in-progress spec-to-pr workflow after a context interruption.

Arguments: $ARGUMENTS
(Format: [JIRA-KEY])
JIRA-KEY is optional — it will be read from .claude/sdd-contract.md if not provided.

---

## Step 1 — Read the Contract Document

Read `.claude/sdd-contract.md` in full. This is the single source of truth for the interrupted session.

If `.claude/sdd-contract.md` does not exist:
- If a JIRA-KEY was provided, offer to restart from Phase 1 (`/spec-to-pr JIRA-KEY`)
- If no JIRA-KEY was provided, tell the user the contract document is missing and ask them to run `/spec-to-pr JIRA-KEY` from scratch

Extract from the contract document:
- JIRA-KEY
- Git stamp (Branch and SHA fields under `## Git stamp`)
- Implementation mode (GREENFIELD / MODIFICATION)
- Layer structure and file paths
- All contract definitions (data layer, logic layer, shared types)
- Codebase conventions (error handling, logging, DB client, test/mock pattern)
- Approved plan for each layer

**Staleness check** — run these two commands and compare against the Git stamp:

```
git branch --show-current   → current branch
git rev-parse --short HEAD  → current SHA
```

| Condition | Action |
|---|---|
| Branch differs from stamp | **STOP** — warn the user before proceeding (see below) |
| Branch matches, SHA differs | Surface an info note — contract may be slightly stale |
| Branch and SHA both match | No warning needed |

**Branch mismatch warning** (halt and display before Step 2):

```
⚠️  Stale contract detected

The contract was written on branch: <stamp branch>
You are currently on branch:        <current branch>

This contract belongs to a different branch. Resuming here will
apply implementation plans intended for a different context.

Options:
  1. Switch to the correct branch: git checkout <stamp branch>
  2. Discard this contract and start fresh: /spec-to-pr <JIRA-KEY>
  3. Proceed anyway (not recommended) — type "force-resume"
```

Wait for user input. Only continue if they type "force-resume".

**SHA mismatch note** (display in Step 3 report, do not halt):

```
ℹ️  Contract written at <stamp SHA>, HEAD is now <current SHA>.
    New commits may have changed files the plan depends on — review the plan before approving.
```

---

## Step 2 — Detect Current Phase

Determine which phase was in progress when the context was lost, by checking for the presence of sentinel files:

| File present | Interpretation |
|---|---|
| No `sdd-contract.md` | Before Phase 4A — nothing implemented yet |
| `sdd-contract.md` exists, no `sdd-data-done.md` or `sdd-logic-done.md` | Phase 4A complete, Phase 4B not started |
| Only `sdd-data-done.md` exists | Phase 4B partial — Data done, Logic still running or not started |
| Only `sdd-logic-done.md` exists | Phase 4B partial — Logic done, Data still running or not started |
| Both `sdd-data-done.md` and `sdd-logic-done.md` exist | Phase 4B complete |
| `sdd-data-gap.md` or `sdd-logic-gap.md` exist | Phase 4B found gaps — resolution needed |
| Transport layer files exist (check paths from contract) | Phase 4C complete |
| Test files exist for all three layers | Phase 4D complete |

Also check the actual implementation files listed in the contract's approved plan — read none of them speculatively, but check for their existence to corroborate the sentinel file evidence.

---

## Step 3 — Report and Confirm Resume Point

Report clearly:

```
## Resuming spec-to-pr — <JIRA-KEY>

**Last completed phase**: <e.g. "Phase 4A — Contracts written">
**Resuming from**: <e.g. "Phase 4B — spawning Data + Logic agents">

**Contract document**: found ✓
**Contract branch**: <stamp branch> <"✓ matches current branch" or "⚠ see staleness note above">
**Contract SHA**: <stamp SHA> <"ℹ N commits behind HEAD" or "✓ up to date">
**Gap files**: <none / list any found>

**Approved plan summary**:
- Data layer: <N files>
- Logic layer: <N files>
- Transport layer: <N files>
- Tests: <N files>
```

Ask: "Resume from Phase <X>? Or correct the detected phase?"

Wait for confirmation before proceeding.

---

## Phase 4B — Data + Logic layers (if resuming here)

Read `.claude/sdd-contract.md` fully. Then spawn two sub-agents simultaneously using the Agent tool.

### Data Agent

Prompt:
```
You are implementing the Data layer for <JIRA-KEY>.

Start by reading `.claude/sdd-contract.md` in full.

Your job:
- Implement every contract method declared in the "Data layer contracts" section
- Write files in the data layer path described in "Layer structure"
- Use the data client and conventions from "Codebase conventions"
- This layer handles IO only — database queries, HTTP calls to external services, file access, etc.
- No business logic here. Translate raw data results to the shared model types from the contract.
- Match existing error handling and logging patterns exactly.

Hard rules:
- Write ONLY to the data layer files. Do not touch contract files, logic layer, or transport layer.
- Do not add methods not declared in the contract.
- If you discover a contract method is missing that you need, write it to `.claude/sdd-data-gap.md` and implement everything else you can.

When done, summarise every file you created or modified in `.claude/sdd-data-done.md`.
```

### Logic Agent

Prompt:
```
You are implementing the Logic layer for <JIRA-KEY>.

Start by reading `.claude/sdd-contract.md` in full.

Your job:
- Implement every contract method declared in the "Logic layer contracts" section
- Write files in the logic layer path described in "Layer structure"
- Business logic, domain rules, and use case orchestration live here
- Call the data layer ONLY via its contract — never import or depend on the concrete data layer implementation
- Use the dependency injection pattern from "Codebase conventions"
- Match existing error handling and logging patterns exactly

Hard rules:
- Write ONLY to the logic layer files. Do not touch contract files, data layer, or transport layer.
- Do not add methods not declared in the contract.
- If you discover a contract method is missing that you need, write it to `.claude/sdd-logic-gap.md` and implement everything else you can.

When done, summarise every file you created or modified in `.claude/sdd-logic-done.md`.
```

After both agents complete, check for gap files. If any exist: resolve gaps in the contract, update `.claude/sdd-contract.md`, then re-run the affected agent.

---

## Phase 4C — Transport layer (if resuming here)

Read `.claude/sdd-contract.md`. Implement the transport layer using the paths from "Approved plan — Transport layer".

Transport rules:
- Parse and validate incoming request → call logic layer via its contract → translate result to response
- No business logic here
- Depend only on the logic layer contract, never the concrete implementation
- Match existing middleware, error translation, and response conventions

---

## Phase 4D — Tests (if resuming here)

Read `.claude/sdd-contract.md`. Spawn three sub-agents simultaneously.

### Data Test Agent
```
You are writing tests for the Data layer of <JIRA-KEY>.
Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.
- Write tests for every method in the data layer
- Use real dependencies (DB, testcontainers, etc.) where the codebase does
- Test: happy path, error cases, edge cases per method
- Target ≥ 80% coverage for the data layer
Write only test files in the data layer path. Do NOT modify any implementation files.
```

### Logic Test Agent
```
You are writing tests for the Logic layer of <JIRA-KEY>.
Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.
- Write tests for every method in the logic layer
- Mock the data layer using its contract — never use the concrete data implementation
- Use the mock/fake convention from "Codebase conventions"
- Test: happy path, business rule enforcement, error propagation
- Target ≥ 80% coverage for the logic layer
Write only test files in the logic layer path. Do NOT modify any implementation files.
```

### Transport Test Agent
```
You are writing tests for the Transport layer of <JIRA-KEY>.
Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.
- Write tests for every endpoint/handler in the transport layer
- Mock the logic layer using its contract
- Test: request parsing, validation rejection, logic error translation, response shape and status codes
- Target ≥ 80% coverage for the transport layer
Write only test files in the transport layer path. Do NOT modify any implementation files.
```

---

## Phase 5 — Quality and PR (if resuming here)

1. Run `sdd mine` to update the knowledge graph
2. Run `sdd quality` — fix every failure before continuing
3. Remove temp files: `.claude/sdd-contract.md`, `.claude/sdd-data-done.md`, `.claude/sdd-logic-done.md`, and any gap files
4. Create a PR with:
   - Title referencing the Jira key
   - Description listing every acceptance criterion and which layer/file satisfies it

---

**Hard rules:**
- Do not re-implement anything already in a `-done.md` file unless its sentinel file is missing or its output files don't exist on disk
- Do not modify the contract document unless resolving a gap file
- If the contract document is corrupted or incomplete, do not guess — tell the user and ask them to re-run `/spec-to-pr JIRA-KEY` from Phase 4A
