# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Run core tests only
go test ./...

# Run all tests (core + all adapters)
make test

# Run adapter tests only
make test-adapter

# Run a single test
go test -run TestGolden ./...

# Update golden YAML fixtures after intentional output changes
make test-update

# Lint (golangci-lint v2.15.0)
make lint
make lint-fix   # auto-fix

# Tidy go.mod for core + all adapters
make tidy

# Full local check: sync + tidy + lint + test
make check

# Install dev tools (golangci-lint)
make install-tools
```

## Architecture

`oaswrap/spec` is a framework-agnostic OpenAPI 3.x document builder for Go. It generates spec documents from route registrations and Go struct reflection rather than parsing code annotations.

### Core package (`github.com/oaswrap/spec`)

- `types.go` — Public interfaces: `Generator` (embeds `Router`), `Router`, `Route`. Also re-exports common OpenAPI types.
- `router.go` — Concrete `generator` struct implementing `Generator`. `NewGenerator`/`NewRouter` are entry points. Routes accumulate in a tree; `build()` is called on every `GenerateSchema`, `MarshalYAML`, `MarshalJSON`, or `Validate` call (not cached — rebuilds each time).
- `errors.go` — `ValidationErrors` aggregating `validate.Error` with severity. Only `SeverityError` items cause `Validate()` / serialization to fail; `SeverityWarning` and `SeverityInfo` are informational.

### Sub-packages

| Package | Role |
|---|---|
| `openapi/` | Data model structs for OpenAPI documents (`Document`, `Schema`, `Operation`, `Config`, etc.). `Config` drives generator behavior. |
| `option/` | Functional options (`OpenAPIOption`, `OperationOption`, `GroupOption`, `ContentOption`) for configuring the generator and individual operations. |
| `internal/builder/` | Converts accumulated route + option data into `openapi.Document` operations via `Builder.AddOperation` and `Builder.AddWebhookOperation`. |
| `internal/reflect/` | Reflects Go types to `openapi.Schema` objects, managing `$components/schemas` de-duplication. |
| `internal/validate/` | Document validation. Issues carry a `Severity` (Error / Warning / Info). `ValidateDocument` is called at the end of every `build()`. |
| `pkg/parser/` | `ColonParamParser` converts `:param` style paths (used by some frameworks) to `{param}` OpenAPI path template format. |
| `internal/testutil/` | Golden-file test helpers used in `_test.go` files. |

### Adapters (`adapter/`)

Eight framework adapters wrap a spec `Generator` alongside a real HTTP router:

```
chiopenapi  echoopenapi  echov5openapi  fiberopenapi
fiberv3openapi  ginopenapi  httpopenapi  muxopenapi
```

Each adapter is a separate Go module (its own `go.mod`). All are included in the root `go.work` workspace. The pattern is:
- `NewRouter(frameworkRouter, ...option.OpenAPIOption)` → returns adapter's `Generator`
- Route registrations go to both the HTTP router and the spec generator simultaneously
- `gen spec.Generator` field must be propagated into every sub-router and group created by the adapter

### Golden tests

Test fixtures live in `testdata/` as `{case}.{version}.yaml` files (e.g. `petstore.v31.yaml`). Tests cover all three version families: `v30`, `v31`, `v32`. Run `make test-update` to regenerate them after intentional output changes.

### Supported OpenAPI versions

- 3.0.x (`openapi.Version304` is the default)
- 3.1.x (`openapi.Version312` recommended)
- 3.2.0 (`openapi.Version320`) — adds `QUERY` method and `$self`

Webhooks require 3.1.x or 3.2.0.

## Key constraints

- Commits must follow [Conventional Commits](https://www.conventionalcommits.org/): `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`, etc. (enforced by lefthook).
- Max line length is 120 characters (golines).
- Import groups: stdlib → third-party → `github.com/oaswrap/spec` (enforced by goimports).
- Release is a two-stage process: `make release-prepare VERSION=x.y.z` (tags core, syncs adapter deps) then `make release-publish VERSION=x.y.z` (tags all adapters).
