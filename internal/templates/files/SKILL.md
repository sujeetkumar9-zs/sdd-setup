# SDD Skill — Spec-Driven Development

You are in SDD mode. Your job: take a Jira issue (plus optional Confluence docs) and deliver a merged PR.

## Rules

- Use `mempalace` MCP for all codebase knowledge — never use Grep, Glob, or Read to explore
- Only read a file when you are about to edit it
- Never write code before presenting a plan and getting explicit approval
- Never add scope beyond the spec
- Match existing patterns exactly — naming, error handling, logging
- If the spec is ambiguous, ask before assuming
- If a framework or pattern is unclear, search mempalace before forming any opinion
- **Never modify shared contracts (interfaces/types) after parallel agents have started** — changes require re-running the affected agent

## Architecture Discovery and Mode Detection

Do NOT assume a folder structure. After gathering context, detect the implementation mode:

**GREENFIELD** — Jira spec describes creating something new AND mempalace finds no existing entity/structure.
**MODIFICATION** — mempalace finds existing layer structure for this entity.

State the mode explicitly before planning.

**If GREENFIELD + architecture doc provided:**
- Use the architecture doc as the authoritative folder/layer blueprint
- Search mempalace ONLY for shared conventions (error handling, logging, DB client, test patterns)
- Show the proposed scaffold in the plan for user confirmation

**If GREENFIELD + no architecture doc:**
- Infer the folder pattern from similar existing entities in mempalace
- Show the inferred scaffold in the plan for user confirmation

**If MODIFICATION:**
- Discover structure entirely from mempalace — paths, naming, patterns already present
- Use architecture doc (if provided) to validate alignment and flag deviations; do not restructure without approval

The multi-agent implementation structure below applies to whatever layers exist — layer names and file paths are always derived from what you discover or the architecture doc, never assumed.

## Workflow

### Phase 1 — Context Gathering (parallel)
- Fetch Jira issue, Confluence pages, and mempalace codebase context simultaneously
- From mempalace: discover the layer structure, contract patterns, error handling, test conventions, and whether the target entity already exists
- After context is gathered: detect and state the implementation mode (GREENFIELD or MODIFICATION)

### Phase 2 — Plan
Build a plan across the discovered layers. Identify:
1. **Contract layer** — where the layer boundaries/contracts are defined (interfaces, types, abstractions). This must be implemented first.
2. **Data layer** — data access, persistence, external IO. Implements the contracts.
3. **Logic layer** — business logic, use cases. Implements the contracts, calls data layer via contract only.
4. **Transport layer** — HTTP/gRPC/CLI handlers. Calls logic layer via contract only.

Include complete contract signatures in the plan (full interface/type/abstract definitions). These are reviewed at approval and drive parallel implementation.

### Phase 3 — Approval Loop
Present the full plan including all contract signatures. Get explicit approval. No code until approved.

### Phase 4A — Contracts (you, sequential — blocks everything)
Implement the contract layer first. Write all interface definitions, shared types, and abstractions that the other layers depend on.

Produce `.claude/sdd-contract.md` containing:
- All contract definitions (verbatim from the written code)
- Layer file paths from the approved plan
- Codebase conventions (DB client, error pattern, test/mock pattern, logging)
- Approved plan excerpts for each remaining layer

**Do not start Phase 4B until the contract document is complete.**

### Phase 4B — Data + Logic layers (two agents, parallel)
Spawn two sub-agents simultaneously, each given `.claude/sdd-contract.md`.

- **Data Agent**: implements the data access contracts. IO only, no business logic.
- **Logic Agent**: implements the business logic contracts. Calls data layer via contract only, never via concrete type.

Each agent writes only to its own layer's files. If an agent discovers a missing contract method, it records it in a gap file rather than inventing one.

Check gap files after both complete. Resolve any gaps in the contract layer, then re-run the affected agent.

### Phase 4C — Transport layer (you, sequential)
Implement handlers/controllers/routes. Wire to logic layer via contract. No business logic here.

### Phase 4D — Tests (three agents, parallel)
Spawn test agents for each layer simultaneously:
- Data layer tests: test against real dependencies (DB, external service) where possible
- Logic layer tests: mock the data layer using the contract
- Transport layer tests: mock the logic layer using the contract

Target: ≥ 80% coverage per layer.

### Phase 5 — Quality and PR
1. `sdd mine` — update knowledge graph
2. `sdd quality` — all gates must pass
3. Remove `.claude/sdd-contract.md` and all temp files
4. Create PR referencing the Jira key and all acceptance criteria
