# sdd — Spec-Driven Development CLI

A CLI tool that sets up AI-powered feature development for Go projects using Claude Code and Mempalace.

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

**From source**
```bash
go install github.com/sujeetkumar9-zs/sdd-setup@latest
```

Verify:
```bash
sdd --version
```

## Requirements

- Go 1.21+
- Python 3.8+
- Claude Code CLI (`npm install -g @anthropic-ai/claude-code`)
- Docker
- Mempalace CLI (`pipx install mempalace`)
- `golangci-lint`, `mockery` (installed automatically by `sdd setup`)

## Usage

### One-time setup

Run once from the root of your Go project:

```bash
cd your-go-project
sdd setup
```

This will:
- Check all prerequisites
- Create a Python virtual environment (`.venv`)
- Install Go quality tools (`golangci-lint`, `mockery`, `goimports`)
- Initialize the Mempalace knowledge graph (`.mempalace/palace/`)
- Mine your codebase into the knowledge graph
- Configure Claude MCP servers
- Install SDD skills and slash commands (`SKILL.md`)
- Update `.gitignore`

To skip the codebase mining step:
```bash
sdd setup --skip-mine
```

### Daily workflow

```bash
# After implementing a feature — update the knowledge graph
sdd mine

# Before creating a PR — run all quality gates
sdd quality

# Check that everything is wired up correctly
sdd verify

# Show knowledge graph and MCP server status
sdd status
```

## Commands

| Command | Description |
|---------|-------------|
| `sdd setup` | One-time project setup |
| `sdd mine` | Re-index codebase into knowledge graph |
| `sdd verify` | Check all SDD components are working |
| `sdd quality` | Run all quality gates before a PR |
| `sdd status` | Show knowledge graph and MCP status |

## Quality Gates

`sdd quality` runs the following checks in order and blocks on any failure:

1. `go test ./...` — all tests must pass
2. `go test -race ./...` — race condition check
3. `go vet ./...` — static analysis
4. `golangci-lint run` — lint checks
5. `gofmt -l .` — formatting check
6. `go build ./...` — must compile

## How It Works

```
Spec (Jira / Confluence)
        ↓
Claude reads spec via Atlassian MCP
        ↓
Claude queries your codebase via Mempalace MCP
        ↓
Claude generates implementation plan → you approve
        ↓
Claude implements code following your patterns
        ↓
Claude generates tests (≥ 80% coverage)
        ↓
Quality gates run automatically
        ↓
PR created and merged ✅
```

## Troubleshooting

Run `sdd verify` to diagnose issues. It checks:

- Mempalace is installed and the palace directory exists
- Palace index has data (`sdd mine` if empty)
- Mempalace MCP server is registered with Claude
- Claude Code is installed
- Go quality tools are available
- Python venv exists
- `SKILL.md` is present
- `.gitignore` is configured

Each failed check includes a specific fix hint.

## Development

```bash
git clone https://github.com/sujeetkumar9-zs/sdd-setup.git
cd sdd-setup
go mod tidy
go build -o sdd .
```
