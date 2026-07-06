Implement the feature described in the following spec from start to PR using the SDD workflow.

Spec: $ARGUMENTS

Steps:
1. Fetch and read the spec (Jira ticket or Confluence page) using the Atlassian MCP
2. Query the codebase using Mempalace MCP to understand existing patterns
3. Generate an implementation plan and present it — wait for approval
4. Implement the feature following existing code conventions
5. Write tests targeting ≥ 80% coverage
6. Run quality gates (go test, go vet, golangci-lint, gofmt)
7. Create a PR with a clear description referencing the spec

Do not proceed past the plan step without explicit approval.
