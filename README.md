# sdd — Spec-Driven Development CLI

A CLI tool that sets up AI-powered feature development for any project using Claude Code and Mempalace. Supports Go, Node/TypeScript, Python, Rust, and Java.

## Install

**macOS (Apple Silicon)**
```bash
curl -fsSL https://github.com/sujeetkumar9-zs/sdd-setup/releases/latest/download/sdd-macos-arm64 \
  -o /usr/local/bin/sdd && chmod +x /usr/local/bin/sdd
```

**macOS (Intel)**
```bash
curl -fsSL https://github.com/sujeetkumar9-zs/sdd-setup/releases/latest/download/sdd-macos-amd64 \
  -o /usr/local/bin/sdd && chmod +x /usr/local/bin/sdd
```

**Linux**
```bash
curl -fsSL https://github.com/sujeetkumar9-zs/sdd-setup/releases/latest/download/sdd-linux-amd64 \
  -o /usr/local/bin/sdd && chmod +x /usr/local/bin/sdd
```

Verify:
```bash
sdd --version
```

## Requirements

- Python 3.8+ (`brew install python` on macOS)
- Claude Code CLI (`npm install -g @anthropic-ai/claude-code`)
- Mempalace CLI (`pipx install mempalace`)
- Atlassian MCP configured (for Jira + Confluence access)

**Language-specific prerequisites** (checked automatically by `sdd setup`):

| Language | Required |
|---|---|
| Go | Go 1.21+ |
| Node / TypeScript | Node + npm |
| Python | Python 3.8+ |
| Rust | Rust toolchain via `rustup` |
| Java | Maven or Gradle |

## Usage

### One-time setup

Run once from the root of your project:

```bash
cd your-project
sdd setup
```

`sdd setup` automatically detects your project language and:

- Checks all prerequisites
- Creates a Python virtual environment (`.venv`)
- Installs language quality tools (see table below)
- Initializes the Mempalace knowledge graph (`.mempalace/palace/`)
- Mines your codebase into the knowledge graph
- Configures Claude MCP servers (mempalace + Atlassian)
- Installs SDD skills and slash commands into `.claude/`
- Updates `.gitignore`

To skip the codebase mining step:
```bash
sdd setup --skip-mine
```

If setup is interrupted (network error, missing tool, etc.), resume from where it left off:
```bash
sdd setup --resume
```

### Daily workflow

```bash
# After implementing a feature — update the knowledge graph
sdd mine

# Mine only changed files (faster on large repos)
sdd mine --changed

# Before creating a PR — run all quality gates
sdd quality

# After upgrading the sdd binary — refresh skills and commands
sdd update

# Check that everything is wired up correctly
sdd verify

# Show knowledge graph and MCP server status
sdd status

# Search the knowledge graph directly
sdd search "repository pattern"

# Remove SDD setup from this project
sdd teardown
```

## CLI Commands

| Command | Description |
|---|---|
| `sdd setup` | One-time project setup |
| `sdd setup --resume` | Resume an interrupted setup |
| `sdd update` | Refresh skills and commands to the latest version |
| `sdd mine` | Re-index codebase into the knowledge graph |
| `sdd mine --changed` | Re-index only git-modified directories (faster) |
| `sdd quality` | Run all quality gates before a PR |
| `sdd verify` | Check all SDD components are working |
| `sdd status` | Show knowledge graph and MCP status |
| `sdd search <query>` | Search the project knowledge graph |
| `sdd teardown` | Remove Mempalace and venv from this project |
| `sdd teardown --full` | Also removes `.claude/commands/` and `SKILL.md` |

## Claude Slash Commands

After `sdd setup`, the following slash commands are available inside Claude Code:

### `/spec-to-pr JIRA-KEY [confluence-page ...]`

Implements a feature end-to-end from a Jira spec. Phases:

1. **Context** — fetches Jira issue, Confluence pages, and queries mempalace for relevant codebase patterns
2. **Plan** — builds a layered implementation plan (models → data → service → API → tests → PR)
3. **Approval** — presents the plan; waits for explicit approval before writing any code
4. **Implementation** — implements each layer following existing codebase conventions
5. **Quality + PR** — runs `sdd quality`, then creates the PR

```
/spec-to-pr PROJ-123
/spec-to-pr PROJ-123 confluence-page-id
/spec-to-pr PROJ-123 arch-page-id api-contracts-page-id
```

### `/spec-to-pr-status`

Reports current progress through the spec-to-pr workflow: which phase is active, what is done, what remains, and any blockers.

### `/spec-to-pr-resume [JIRA-KEY]`

Recovers a `spec-to-pr` session that was interrupted by a context window overflow. Reads `.claude/sdd-contract.md` (the handoff artifact written at the end of Phase 4A) to determine what was already completed, then resumes from the correct phase without re-implementing work that is already done.

- If `sdd-contract.md` is missing and a JIRA-KEY is provided, offers to restart from Phase 1
- Detects phase by checking sentinel files (`sdd-data-done.md`, `sdd-logic-done.md`, gap files) and actual file existence on disk
- Confirms the detected resume point with the user before spawning any agents
- Respects the original approved plan — does not re-plan or re-prompt for approval

```
/spec-to-pr-resume
/spec-to-pr-resume PROJ-123
```

### `/spec-to-pr-quality`

Runs all language-appropriate quality gates and reports results with fix instructions for each failure.

### `/fix-bug JIRA-KEY [confluence-page ...]`

Investigates and fixes the bug described in a Jira ticket. Phases:

1. **Context** — fetches the Jira bug report (description, repro steps, stack traces), Confluence architecture docs, and queries mempalace to map the affected code path through all layers
2. **Debug** — traces the reported behaviour through the code path; identifies the defect location and logic flaw using mempalace; forms a ranked list of root cause hypotheses
3. **RCA** — presents a structured Root Cause Analysis (entry point, defect location, exact logic flaw, evidence, impact); **waits for user confirmation before planning any fix**
4. **Fix plan** — the minimum correct change, file by file, with a regression test covering the exact repro scenario; waits for explicit approval before writing code
5. **Implementation** — makes exactly the approved changes, writes the regression test, runs `sdd quality`
6. **PR** — title `fix(JIRA-KEY): <summary>`, description includes RCA, what changed, regression test, and manual verification steps

```
/fix-bug PROJ-123
/fix-bug PROJ-123 confluence-page-id
```

### `/address-pr-comments [JIRA-KEY] [confluence-page ...]`

Addresses all unresolved reviewer comments on the open PR for the current branch. Phases:

1. **Context** — fetches the PR diff and all unresolved review comments, loads comment history (`.claude/pr-feedback-history.md`) to skip already-addressed items, queries mempalace for patterns in touched files
2. **Triage** — classifies every comment as ACTIONABLE, AMBIGUOUS, or OUT OF SCOPE; batches all clarifying questions in one shot before planning; surfaces out-of-scope requests for user decision
3. **Plan** — one entry per comment: exact file, line range, and change description; presents a summary table and waits for explicit approval before writing any code
4. **Implementation** — makes only the changes in the approved plan; runs `sdd quality` and fixes any regressions
5. **History** — appends a session record to `.claude/pr-feedback-history.md` (addressed / deferred / out of scope) so future runs skip what was already done

```
/address-pr-comments
/address-pr-comments PROJ-123
/address-pr-comments PROJ-123 confluence-page-id
```

### `/review-pr JIRA-KEY [confluence-page ...]`

Reviews the open PR on the current branch as a senior engineer. Phases:

1. **Context** — fetches the PR diff, Jira acceptance criteria, Confluence architecture docs, and queries mempalace for patterns in the touched layers. Reads any external repos or frameworks referenced in the spec.
2. **Review** — evaluates the diff across eight dimensions:
   - **Requirement fulfillment** — every acceptance criterion traced to the diff
   - **Architecture alignment** — compliance with existing patterns and any provided architecture docs
   - **SOLID principles** — SRP, OCP, LSP, ISP, DIP where applicable
   - **Design patterns** — correct use of patterns; flags anti-patterns and suggests better alternatives
   - **Code quality** — naming, error handling, logging, dead code
   - **Test coverage** — happy path, error paths, edge cases, meaningful assertions
   - **Security** — input validation, injection risks, auth checks, secret handling
   - **Performance** — N+1 queries, unbounded loops, missing timeouts
3. **Report** — produces a structured review with verdict (`APPROVE` / `REQUEST CHANGES` / `NEEDS DISCUSSION`), file:line-anchored findings, corrected code snippets, and explicit callout of what was done well
4. **Confirmation** — presents the review in the conversation; only posts to GitHub if the user explicitly confirms

```
/review-pr PROJ-123
/review-pr PROJ-123 confluence-page-id
/review-pr PROJ-123 arch-page-id api-contracts-page-id
```

## Quality Gates

`sdd quality` (and `/spec-to-pr-quality`) runs the appropriate gates for the detected language:

**Go**
1. `go test ./...` — all tests must pass
2. `go test -race ./...` — race condition check
3. `go vet ./...` — static analysis
4. `golangci-lint run` — lint checks
5. `gofmt -l .` — formatting check
6. `go build ./...` — must compile

**Node / TypeScript**
1. `npm test` — all tests must pass
2. `npm run lint` — lint checks
3. `npm run build` — must compile/bundle

**Python**
1. `pytest` — all tests must pass
2. `ruff check .` — lint checks
3. `ruff format --check .` — formatting check

**Rust**
1. `cargo test` — all tests must pass
2. `cargo clippy -- -D warnings` — lint checks
3. `cargo fmt --check` — formatting check
4. `cargo build` — must compile

**Java (Maven)**
1. `mvn test` — all tests must pass
2. `mvn package -DskipTests` — must build

**Java (Gradle)**
1. `./gradlew test` — all tests must pass
2. `./gradlew build -x test` — must build

## Language Quality Tools

Installed automatically by `sdd setup`:

| Language | Tools |
|---|---|
| Go | `golangci-lint`, `mockery`, `goimports` |
| Node / TypeScript | `eslint` |
| Python | `ruff`, `pytest` |
| Rust | uses `rustup` (no additional install) |
| Java | uses Maven or Gradle (no additional install) |

## How It Works

### Spec-to-PR workflow

```
Jira issue + Confluence pages
        ↓
Claude fetches spec via Atlassian MCP
        ↓
Claude queries codebase via Mempalace MCP
        ↓
Layered implementation plan → you approve
        ↓
Claude implements code matching your patterns
        ↓
Tests written (≥ 80% coverage target)
        ↓
Quality gates run automatically
        ↓
PR created ✅
```

### Review-PR workflow

```
Jira issue + Confluence pages + current branch PR
        ↓
Claude fetches PR diff via GitHub MCP
        ↓
Claude fetches spec + architecture docs via Atlassian MCP
        ↓
Claude queries codebase context via Mempalace MCP
        ↓
8-dimension review (requirements, architecture, SOLID,
patterns, quality, tests, security, performance)
        ↓
Structured review report with verdict + fix suggestions
        ↓
You confirm → posted to GitHub as PR review ✅
```

## Troubleshooting

Run `sdd verify` to diagnose issues. It checks:

- Mempalace is installed and the palace directory exists
- Palace index has data (`sdd mine` if empty)
- Mempalace MCP server is registered with Claude
- Claude Code is installed
- Language quality tools are available
- Python venv exists
- SDD skills and commands are present in `.claude/`
- `.gitignore` is configured

Each failed check includes a specific fix hint.

## Development

> **Note**: Go is only required if you are building `sdd` from source. End users install the pre-compiled binary via `curl` above and do not need Go.

```bash
git clone https://github.com/sujeetkumar9-zs/sdd-setup.git
cd sdd-setup
go mod tidy
go build -o sdd .
./sdd --help
```

To test changes in another project without publishing a release:

```bash
go build -o sdd .   # rebuild binary
cd your-other-project
/path/to/sdd setup  # use the local binary directly
```
