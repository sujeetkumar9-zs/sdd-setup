Implement the feature described in the spec using the SDD workflow.

Arguments: $ARGUMENTS
(Format: JIRA-KEY [confluence-page-id-or-url ...])

---

## Phase 1 — Context Gathering

Run these three in parallel:

**Jira**: Fetch the issue (first argument). Extract:
- Summary, description, acceptance criteria
- Labels, priority, linked issues, referenced component or service names

**Confluence**: For each Confluence argument (all args after the Jira key):
- Fetch the page. If empty or no useful content → fetch its child pages instead.
- Extract: data models, service decisions, architecture patterns, API contracts, field definitions
- If an architecture or design page is provided, treat it as the authoritative guide for layer structure

**Codebase** (mempalace):

Before querying mempalace, check for reusable context:

1. **If `.claude/sdd-contract.md` exists** — read the `## Codebase conventions` and `## Layer structure` sections. Use those values directly. Skip convention queries below.
2. **If codebase conventions were already gathered earlier in this conversation** (layer structure, error handling, logging, test/mock pattern were described in a previous skill or query result) — skip convention queries below.

**Spec-specific queries** (always run):
- Types, interfaces, and services mentioned in the spec
- Whether the entity/service named in the spec already exists in the codebase

**Convention queries** (skip if reusable context found above):
- The project's layer structure — how it separates data access, business logic, and transport/API concerns
- The naming convention for each layer (e.g. `repository/`, `service/`, `controller/` — or `store/`, `core/`, `handler/`)
- How contracts between layers are expressed: Go interfaces, TypeScript types/interfaces, Python abstract base classes, Java interfaces, etc.
- Existing error handling, logging, and test/mock conventions

If any framework, library, or pattern is unclear: search mempalace before forming any opinion. Do NOT guess.

**Do not assume a folder structure. Derive everything from what you find.**

---

## Spec Quality Gate (run after Phase 1 — before planning anything)

Evaluate the Jira ticket against four dimensions. Check each one explicitly.

**Dimension 1 — Measurable acceptance criteria**
Every AC must be a verifiable condition with a clear pass/fail. Flag any AC that uses vague language without a specific target:
- "improve", "make better", "more reliable", "cleaner", "faster" — without a metric (e.g. p99 < 200ms, error rate < 0.1%)
- "users should be able to…" without specifying the exact behaviour
- "handle edge cases" without naming them

**Dimension 2 — No deferred knowledge**
Flag any of the following in the description or AC:
- "ask [person/team]", "check with [name]", "talk to [team]"
- "TBD", "to be defined", "to be discussed", "details pending"
- "see [person] for details", "per [name]'s decision"
- "design doc forthcoming", references to meetings or Slack threads

**Dimension 3 — Defined scope**
Flag if:
- The description is fewer than 2 sentences with no acceptance criteria
- It is unclear whether this is a new feature, a change to existing behaviour, or a bug fix
- "maybe", "possibly", "we could also" appears in the scope — undefined optionality

**Dimension 4 — Technical context**
Flag if the ticket provides no pointer to affected code or system — no service name, no API endpoint, no database entity, no Confluence page link, and none were passed as arguments. (Mempalace context alone is sufficient if Phase 1 found the entity.)

---

**If all four dimensions pass:** proceed to Mode Detection.

**If any dimension fails:** halt and report:

```
## Spec Quality Gate — FAILED

The Jira ticket has gaps that increase the risk of building the wrong thing.

### Gaps found
- [Dimension N] <specific problem>
- ...

### What to do
1. Update the Jira ticket with the missing information
2. Re-run: /spec-to-pr <JIRA-KEY>

Or, if you want to proceed despite these gaps:
Type "override" to continue — you accept the risk of misaligned implementation.
```

Wait for user input. If they type "override", continue to Mode Detection and note the gaps in the plan. Otherwise, stop.

---

## Mode Detection (after Phase 1 context gathering)

Before planning, determine the implementation mode. State the detected mode explicitly.

**GREENFIELD** — if ALL of the following are true:
- The Jira spec describes creating a new service, entity, or module (language like "create", "implement new", "scaffold", "build a new")
- Mempalace finds no existing folder structure, types, or services for this entity

**MODIFICATION** — if ANY of the following are true:
- Mempalace finds existing layer structure, types, interfaces, or services for this entity
- The Jira spec describes changing, extending, or fixing existing behaviour

State the mode clearly before Phase 2, e.g.:
> **Mode: GREENFIELD** — no existing entity found in mempalace; Jira spec says "create new X service"

or:

> **Mode: MODIFICATION** — mempalace found existing layer structure at `entities/orders/`

---

## Phase 2 — Plan

Based on what you discovered **and the detected mode**, identify the layer boundaries:

1. **Contract layer** — where shared types and layer-boundary contracts are defined. This must be written first because everything else depends on it.
2. **Data layer** — data access, persistence, external IO. Implements the data contracts.
3. **Logic layer** — business logic and use cases. Implements the logic contracts, calls data layer via contract.
4. **Transport layer** — HTTP/gRPC/CLI/etc. Calls logic layer via contract.

**Planning behavior depends on mode:**

### If GREENFIELD + architecture Confluence page provided
- Use the architecture doc as the authoritative blueprint for folder structure, layer naming, and file naming conventions
- Do NOT try to discover layer structure from mempalace — the entity does not exist yet
- Search mempalace ONLY for conventions shared across existing entities: error handling pattern, logging calls, DB/HTTP client injection, test/mock approach
- Show the proposed folder scaffold explicitly in the plan so the user can correct it before approving

### If GREENFIELD + no architecture doc provided
- Search mempalace for the dominant folder/layer pattern used by other entities in the codebase
- Use that pattern as the template for the new entity
- Show the inferred scaffold in the plan for user confirmation

### If MODIFICATION
- Discover the existing structure entirely from mempalace — use paths, naming, and patterns already there
- If an architecture doc was provided, use it to validate alignment and flag any deviations in the plan (do not restructure without user approval)
- Show only the files to be added or changed — do not re-scaffold the whole entity

---

For each layer, state:
- Exact file paths to create or modify (using the naming this codebase already uses, or the architecture doc if GREENFIELD)
- Complete contract signatures — every interface method, type, or abstract definition that other layers depend on
- Rationale tied to a specific acceptance criterion

**Include the full contract signatures in the plan.** These are what the user approves and what parallel agents implement against. If they are incomplete here, agents will collide mid-implementation.

Present the complete plan.

---

## Phase 3 — Approval Loop

1. Present the complete plan including all contract signatures.
2. Ask: "Approve this plan, suggest changes, or ask questions?"
3. If changes → revise (pay close attention to contract completeness) and re-present.
4. Repeat until the user gives explicit approval.
5. **Do NOT write any code until approved.**

---

## Phase 4A — Contracts (You implement this directly)

Implement the contract layer first. This is the shared foundation everything else builds on.

Write all type definitions, interfaces, and abstractions that define the boundaries between layers — in whatever form this language and codebase use. Match existing conventions exactly.

After the contract files are written, run these two commands to capture the git stamp:
- `git branch --show-current` → current branch name
- `git rev-parse --short HEAD` → current commit SHA

Then produce `.claude/sdd-contract.md`:

```
# SDD Contract — <JIRA-KEY>

## Git stamp
Branch: <output of git branch --show-current>
SHA: <output of git rev-parse --short HEAD>

## Implementation mode
<GREENFIELD or MODIFICATION — and one sentence explaining why>

## Layer structure
<If GREENFIELD: "Scaffolded from architecture doc at <page-id> / Inferred from codebase pattern". If MODIFICATION: "Discovered from mempalace — existing entity at <path>">
<describe the layer naming and folder structure in use>

## Language and contract style
<e.g. "Go interfaces", "TypeScript interface types", "Python ABCs", "Java interfaces">

## Data layer contracts
<paste the full contract definitions for the data layer verbatim>

## Logic layer contracts
<paste the full contract definitions for the logic layer verbatim>

## Shared types / models
<list all shared types with their fields/signatures>

## Codebase conventions
- Error handling: <pattern found in codebase>
- Logging: <logger and call pattern>
- Data client: <DB/HTTP client used and how it is injected>
- Test/mock pattern: <how mocks or fakes are structured in this codebase>
- Dependency injection: <constructor pattern, DI framework, or wire-up convention>

## Approved plan — Data layer
<copy verbatim from the approved plan>

## Approved plan — Logic layer
<copy verbatim from the approved plan>

## Approved plan — Transport layer
<copy verbatim from the approved plan>

## Approved plan — Tests
<copy verbatim from the approved plan>
```

**Do not proceed to Phase 4B until `.claude/sdd-contract.md` exists and is complete.**

---

## Phase 4B — Data + Logic layers (Spawn two sub-agents in parallel)

Read `.claude/sdd-contract.md` fully. Then spawn two sub-agents simultaneously using the Agent tool.

### Data Agent

Prompt:
```
You are implementing the Data layer for <JIRA-KEY>.

Start by reading `.claude/sdd-contract.md` in full.

Your job:
- Implement every contract method declared in the "Data layer contracts" section
- Write files in the data layer path described in "Layer structure discovered"
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
- Write files in the logic layer path described in "Layer structure discovered"
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

After both agents complete, check for gap files (`.claude/sdd-data-gap.md`, `.claude/sdd-logic-gap.md`). If any exist:
1. Add the missing methods to the appropriate contract files
2. Re-run the affected agent with the updated contract

---

## Phase 4C — Transport layer (You implement this directly)

Read `.claude/sdd-contract.md`. Implement the transport layer (HTTP handlers, gRPC handlers, CLI commands, etc.) using the paths from "Approved plan — Transport layer".

Transport rules:
- Parse and validate incoming request → call logic layer via its contract → translate result to the appropriate response format
- No business logic here — all decisions belong in the logic layer
- Depend only on the logic layer contract, never on the concrete implementation
- Match existing middleware, error translation, and response conventions

---

## Phase 4D — Tests (Spawn three sub-agents in parallel)

Read `.claude/sdd-contract.md`. Spawn three sub-agents simultaneously.

### Data Test Agent

Prompt:
```
You are writing tests for the Data layer of <JIRA-KEY>.

Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.

Your job:
- Write tests for every method in the data layer
- Use the test/mock pattern from "Codebase conventions" — use real dependencies (DB, testcontainers, etc.) where the codebase does
- Test: happy path, error cases, edge cases per method
- Target ≥ 80% coverage for the data layer

Write only test files in the data layer path. Do NOT modify any implementation files.
```

### Logic Test Agent

Prompt:
```
You are writing tests for the Logic layer of <JIRA-KEY>.

Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.

Your job:
- Write tests for every method in the logic layer
- Mock the data layer using its contract from the contract document — never use the concrete data implementation
- Use the mock/fake convention from "Codebase conventions"
- Test: happy path, business rule enforcement, error propagation
- Target ≥ 80% coverage for the logic layer

Write only test files in the logic layer path. Do NOT modify any implementation files.
```

### Transport Test Agent

Prompt:
```
You are writing tests for the Transport layer of <JIRA-KEY>.

Read `.claude/sdd-contract.md` for contracts, file paths, and test conventions.

Your job:
- Write tests for every endpoint/handler in the transport layer
- Mock the logic layer using its contract from the contract document
- Test: correct request parsing, validation rejection, logic error translation, response shape and status codes
- Target ≥ 80% coverage for the transport layer

Write only test files in the transport layer path. Do NOT modify any implementation files.
```

---

## Phase 5 — Quality and PR

1. Run `sdd mine` to update the knowledge graph
2. Run `sdd quality` — fix every failure before continuing
3. Remove temp files: `.claude/sdd-contract.md`, `.claude/sdd-data-done.md`, `.claude/sdd-logic-done.md`, and any gap files
4. Create a PR with:
   - Title referencing the Jira key
   - Description listing every acceptance criterion and which layer/file satisfies it

---

**Hard rules (apply at every phase):**
- Do not proceed past Phase 3 without explicit user approval
- Do not read any file unless you are about to edit it
- Do not add scope beyond the spec
- Do not modify contracts after Phase 4B starts — changes require re-running the affected agent
