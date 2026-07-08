Implement the feature described in the following spec from start to PR using the SDD workflow.

Spec: $ARGUMENTS

Steps:
1. Fetch and read the spec (Jira ticket or Confluence page) using the Atlassian MCP
2. Search the codebase using the `mempalace` MCP search tool — do NOT use Grep, Glob, Read, or Bash to explore the codebase
3. Generate an implementation plan and present it — wait for explicit approval before writing any code
4. Implement the feature following existing code conventions exactly
5. Write tests targeting ≥ 80% coverage
6. Run `sdd mine` to update the knowledge graph with your new code
7. Run quality gates: go test ./..., go vet ./..., golangci-lint run, gofmt -l .
8. Create a PR with a clear description referencing the spec

Do not proceed past the plan step without explicit approval.
Do not open any file for reading unless you are about to edit that specific file.
