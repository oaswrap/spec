# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Development Tools (install once)
```bash
make install-tools   # installs gotestsum and golangci-lint
```

### Testing
```bash
make test                    # run core + all adapter tests
make test-adapter            # run adapter tests only
make test-update             # update golden files (testdata/*.yaml)
make testcov                 # run tests with coverage
make testcov-html            # generate + open HTML coverage report

# Run a single test (core)
go test ./... -run TestName

# Run a single adapter test
cd adapter/fiberopenapi && go test ./... -run TestName
```

### Code Quality
```bash
make lint    # lint core and all adapters
make tidy    # go mod tidy for core + all adapters
make sync    # go work sync
make check   # sync + tidy + lint + test (full local CI)
```

### Release (two-stage flow)
```bash
make release-prepare VERSION=x.y.z   # Stage 1: tag core + sync adapter go.mod
# then commit and push the adapter go.mod changes
make release-publish VERSION=x.y.z   # Stage 2: tag and push all adapters
make release-dry-run VERSION=x.y.z   # preview without making changes
```

## Architecture

This is a **Go workspace monorepo** (`go.work`) containing:
- **Core module** (`github.com/oaswrap/spec`) — the OpenAPI spec builder
- **Adapter modules** (`adapter/<name>/`) — framework-specific wrappers, each a separate Go module with its own `go.mod`

### Core module

| File/Package | Purpose |
|---|---|
| `types.go` | Public interfaces: `Generator`, `Router`, `Route`, `reflector`, `spec` |
| `router.go` | `generator` struct — implements route registration, grouping, and spec serialization |
| `reflector.go` | Dispatches to OAS 3.0 or 3.1 reflector based on version string |
| `reflector3.go` / `reflector31.go` | Version-specific reflectors wrapping `swaggest/openapi-go` |
| `operation.go` | Builds operation context from options |
| `jsonschema.go` | JSON Schema helpers |
| `option/` | All `option.OpenAPIOption`, `option.OperationOption`, `option.GroupOption` constructors |
| `openapi/` | Internal config structs and OpenAPI entity types |
| `pkg/` | Shared utilities: `mapper`, `parser`, `dto`, `testutil`, `util` |
| `internal/` | Private helpers: `debuglog`, `errs`, `mapper`, `parser` |

### How spec generation works

1. `NewRouter`/`NewGenerator` creates a `generator` with a `reflector` (OAS 3.0 or 3.1).
2. Route methods (`Get`, `Post`, etc.) register `route` objects (lazy — no reflector call yet).
3. `Group` creates child `generator`s sharing the same `reflector`; group options are propagated at build time.
4. On first `MarshalYAML`/`MarshalJSON`/`Validate`, `buildOnce()` fires: flattens the route tree and calls `reflector.Add()` for each route, applying group options.
5. `reflector.Add()` translates `option.OperationOption` functions into the `swaggest/openapi-go` operation model.

### Adapters

Each adapter (e.g. `adapter/fiberopenapi`) wraps both the framework router and `spec.Generator`:
- `NewRouter(frameworkRouter, opts...)` → returns the adapter's `Generator` interface
- Route methods register handlers on the framework router **and** call `spec.Router.Add()` for documentation
- Docs UI (`/docs`) and YAML spec (`/docs/openapi.yaml`) routes are added automatically unless `option.DisableDocs()` is set
- Each adapter uses a `parser.ColonParamParser` (or similar) to convert framework-specific path params (e.g. `:id`) to OpenAPI style (`{id}`)

### Golden-file testing

Tests in `router_test.go` and adapter `router_test.go` compare generated YAML output against files in `testdata/`. Run `make test-update` to regenerate them after intentional changes.

### Versioning & tagging

- Core tag: `v1.2.3`
- Adapter tags: `adapter/<name>/v1.2.3`
- After tagging core, run `make sync-adapter-deps VERSION=x.y.z` to pin adapter `go.mod` to the new core version, then commit before publishing adapter tags.

### Git hooks (lefthook)

Pre-commit: `gofmt`, `go vet`, `golangci-lint`, `go mod tidy`.  
Commit messages must follow Conventional Commits: `feat|fix|docs|style|refactor|test|chore`.
