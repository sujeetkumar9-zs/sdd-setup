# Claude Code Instructions

## Codebase Knowledge — Use Mempalace First

This project has a **Mempalace** MCP server connected (`mempalace`). It contains a pre-indexed semantic knowledge graph of the entire codebase.

### Hard rules

- **NEVER** use `Explore`, `Grep`, `Glob`, `Read`, or `Bash` (find/grep/cat) to understand the codebase structure, find types, or discover patterns
- **ALWAYS** call the `mempalace` MCP search tool first for any question about the codebase
- Only open a specific file with `Read` when you are about to make a targeted edit to that exact file

### When to use Mempalace vs file tools

| Task | Tool to use |
|---|---|
| Understanding a domain concept or flow | `mempalace` MCP search |
| Finding which package/file handles X | `mempalace` MCP search |
| Understanding existing types or interfaces | `mempalace` MCP search |
| Discovering conventions or patterns | `mempalace` MCP search |
| Making an edit to a specific known file | `Read` then `Edit` |

### Example

If asked to understand the household cards endpoint flow:
1. Call `mempalace` MCP → search "household cards endpoint flow"
2. Use those results to form a plan
3. Only `Read` a file if you need to edit it

Do NOT run `grep -r FetchHouseholdCards` or open `service_interfaces.go` to explore — that information is already in Mempalace.

## Go Project Conventions

- Follow existing patterns exactly — naming, error handling, logging style
- Do not introduce new conventions or abstractions not already present
- Match the layer structure: handler → service → repository (or equivalent in this repo)

## Workflow

1. Search Mempalace to understand the relevant domain and patterns
2. Plan changes (list files to modify, types to add)
3. Get user approval before writing any code
4. Implement, then run `go build ./...`, `go test ./...`, `golangci-lint run`
