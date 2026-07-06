Run all quality gates for the current Go project and report results.

Run in order:
1. `go test ./...` — all tests must pass
2. `go test -race ./...` — no race conditions
3. `go vet ./...` — no static analysis issues
4. `golangci-lint run` — no lint errors
5. `gofmt -l .` — no unformatted files (fix with `gofmt -w .`)
6. `go build ./...` — must compile cleanly

For each failure, explain what is wrong and provide the exact command to fix it.
Only confirm ready-for-PR when all six gates pass.
