# BAD_SPEC Example

This is an example of a vague spec that will cause problems in SDD. Do not start implementation until these gaps are resolved.

---

## Jira: PLAT-5678 — Improve the search

### Description

Search is slow and users are complaining. Make it faster and better.
Also add some caching maybe. Talk to John about the details.

---

## Why this spec is insufficient

| Problem | Impact |
|---------|--------|
| No definition of "slow" or "better" | Cannot write acceptance tests |
| "Maybe" caching — scope is undefined | Risk of building the wrong thing |
| "Talk to John" — key decisions are undocumented | Blocks implementation, knowledge is not captured |
| No success criteria | No way to know when the ticket is done |
| No technical constraints or pointers | Must guess at which code to change |

## What to do

Before starting implementation, go back to the author and get:

1. Specific performance targets (e.g. p99 latency < 200ms)
2. A decision on caching: yes or no, and if yes — what layer, what TTL, what invalidation strategy
3. Written acceptance criteria with measurable conditions
4. Pointers to the relevant code or systems
