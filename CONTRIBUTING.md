# Contributing to oaswrap/spec

Thank you for your interest in contributing! This guide covers everything you need to get started.

## Table of Contents

- [Ways to Contribute](#ways-to-contribute)
- [Getting Started](#getting-started)
- [Development Workflow](#development-workflow)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Code Quality](#code-quality)
- [Submitting a Pull Request](#submitting-a-pull-request)
- [Adding a New Adapter](#adding-an-adapter)
- [Commit Messages](#commit-messages)

---

## Ways to Contribute

- **Report bugs** — Open an issue with clear reproduction steps and expected vs actual behavior.
- **Suggest features** — Open an issue describing the use case before writing code.
- **Improve docs** — Fix typos, add examples, or clarify existing documentation.
- **Submit PRs** — Pick up an existing issue or fix a bug you discovered.

Please search existing issues before opening a new one.

For community expectations, see [CODE_OF_CONDUCT.md](./CODE_OF_CONDUCT.md).
For security vulnerabilities, follow [SECURITY.md](./SECURITY.md) and use private disclosure (GitHub Security Advisories) instead of public issues.

---

## Getting Started

### Prerequisites

- Go 1.22 or later
- Git

### Setup

```bash
# Clone the repository
git clone https://github.com/oaswrap/spec.git
cd spec

# Install development tools (golangci-lint)
make install-tools
```

### Verify your setup

```bash
make check   # sync + tidy + lint + test
```

---

## Development Workflow

The repository is a **Go workspace monorepo**. The `go.work` file ties together the core module and all adapter modules.

```
spec/                  # core module (github.com/oaswrap/spec)
adapter/
  chiopenapi/          # Chi adapter
  echoopenapi/         # Echo adapter
  fiberopenapi/        # Fiber adapter
  ginopenapi/          # Gin adapter
  httpopenapi/         # net/http adapter
  muxopenapi/          # Gorilla Mux adapter
  ...
```

Changes to the core module may require corresponding updates to affected adapters. Changes to a single adapter are self-contained.

---

## Testing

```bash
# Run all tests (core + adapters with go test)
make test

# Run adapter tests only
make test-adapter

# Run a single core test
go test ./... -run TestName

# Run a single adapter test
cd adapter/fiberopenapi && go test ./... -run TestName

# Run with coverage
make testcov

# Open HTML coverage report
make testcov-html
```

### Golden files

Many tests compare generated YAML output against files in `testdata/`. If your change intentionally modifies the generated output, regenerate the golden files:

```bash
make test-update
```

Always review the diff of regenerated golden files to confirm the changes are expected.

---

## Code Quality

```bash
make lint    # run golangci-lint on core and all adapters
make tidy    # go mod tidy for core + all adapters
make sync    # go work sync
make check   # run all of the above + tests
```

The pre-commit hook (via [lefthook](https://github.com/evilmartians/lefthook)) runs `gofmt`, `go vet`, `golangci-lint`, and `go mod tidy` automatically. Make sure your code passes all of these before pushing.

---

## Submitting a Pull Request

1. **Fork** the repository and create a feature branch from `main`.
2. **Make your changes**, keeping commits focused and atomic.
3. **Add or update tests** to cover your change. PRs without tests may be asked to add them.
4. **Update golden files** if you changed spec generation output (`make test-update`).
5. **Run the full check** locally: `make check`.
6. **Open a PR** against `main` with a clear description of what changed and why.

PRs must pass CI before merging: core tests on Go 1.22–1.26 plus lint, each adapter's lint and tests on the Go version declared in its `go.mod`, and `govulncheck`/`gosec` scans.

---

## Adding an Adapter

To add support for a new Go web framework:

1. Create `adapter/<frameworkname>openapi/` with its own `go.mod`.
2. Add the new module to `go.work`.
3. Implement the `spec.Generator` interface by wrapping both the framework router and the core `spec.Router`.
4. Register routes on the framework router **and** call `spec.Router.Add()` for documentation.
5. Automatically mount `/docs` (UI) and `/docs/openapi.yaml` (spec) unless `option.WithDisableDocs()` is set.
6. Use a `parser.ColonParamParser` (or equivalent) to translate framework path params (e.g. `:id`) to OpenAPI style (`{id}`).
7. Add a `testdata/` directory and golden-file tests following the pattern in existing adapters.
8. Add a `README.md` describing the adapter.

Look at `adapter/fiberopenapi` or `adapter/ginopenapi` as reference implementations.

---

## Commit Messages

This repository follows [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>: <short description>
```

Allowed types: `feat`, `fix`, `docs`, `style`, `refactor`, `test`, `chore`.

Examples:

```
feat: add response header option
fix: handle empty path group correctly
docs: clarify WithSecurity usage in README
test: add golden file for nested groups
chore: sync adapter deps to v0.5.0
```

The commit-msg hook enforces this format. Breaking changes should include `BREAKING CHANGE:` in the commit body.
