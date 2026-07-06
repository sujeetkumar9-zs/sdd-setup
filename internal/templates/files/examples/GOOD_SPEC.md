# GOOD_SPEC Example

This is an example of a well-written spec that SDD can act on directly.

---

## Jira: PLAT-1234 — Add rate limiting to the product search API

### Background

The product search endpoint (`GET /v1/products/search`) has no rate limiting.
Under load testing it saturates the downstream inventory service at ~500 RPS.
We need per-client rate limiting enforced at the API gateway layer.

### Acceptance Criteria

- [ ] Each API client (identified by `X-Client-ID` header) is limited to 100 requests per second
- [ ] Requests exceeding the limit receive `429 Too Many Requests` with a `Retry-After` header
- [ ] Rate limit counters reset every second (sliding window)
- [ ] Limits are configurable per environment via `RATE_LIMIT_RPS` env var (default: 100)
- [ ] Metrics are emitted: `api.rate_limit.allowed` and `api.rate_limit.rejected` (counter, labelled by client ID)

### Out of Scope

- Per-endpoint limits (all endpoints share the same limit for now)
- Persistent storage of counters across restarts (in-memory is fine)

### Technical Notes

- Use the existing `middleware` package under `internal/middleware/`
- Existing middleware is registered in `internal/server/server.go` — follow the same pattern
- The `X-Client-ID` header is already validated upstream; it will always be present

### Definition of Done

- Unit tests for the middleware covering: allowed request, rejected request, counter reset
- Integration test hitting the endpoint at 2× the limit and asserting 429s
- PR linked to this ticket
