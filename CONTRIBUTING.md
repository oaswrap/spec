# Contributing to oaswrap/spec

Thank you for contributing! This guide explains how to set up your environment, run checks, and safely release new versions.

## 🚀 Quick Start

```bash
# Install required tools
make install-tools

# Run all local quality checks (lint, tidy, tests)
make check

# Run tests in parallel (faster!)
make test-parallel

# Run coverage
make testcov

# Tidy modules
make tidy

# Sync workspace
make sync

# Fix accidental replace statements
make fix-replace
```

## ✅ How We Organize Modules

- **Core** module: root package (`./...`)
- **Adapters**: `adapters/fiberopenapi`, `adapters/ginopenapi`, `adapters/echoopenapi`, `adapters/chiopenapi`

Each adapter has its own `go.mod`. Adapters must **not** have `replace` statements in their `go.mod` — this is enforced by `make check-replace-strict`.

All version links between core and adapters are handled through the workspace (`go.work`).

## 🧹 Required Quality Gates

| Command | What it does |
|----------------|----------------------------------------------|
| `make check` | Lint, tidy, and test all modules |
| `make check-release` | Same as `check` but also drops `replace` statements and enforces strict checks |
| `make release-check` | Verifies your Git state is clean (no uncommitted or staged changes) |

## 🔖 Bump Adapter Versions

When updating the core version (for example, `v0.2.0-dev.1`):

```bash
make bump-dev NEXT=v0.2.0-dev.1
```

This updates all adapters to use the new version.

## 🚢 Tag a Release

**Releases are triggered by pushing a Git tag.**

### ✅ Production release

```bash
# Make sure your working tree is clean!
make release-check

# Run the full release check (incl. fix replaces)
make check-release

# Tag and push:
git tag v1.2.3
git push origin v1.2.3
```

Or use the Makefile helper:

```bash
# Example for dev pre-release
make release-dev VERSION=v0.2.0-dev.1
```

This runs checks, creates the tag, and pushes.  
Pushing the tag automatically triggers `./github/workflows/release.yml` — which runs all quality gates again and creates a GitHub Release with a changelog.

## 🏷️ Dev vs Production

- **Production release**: `v1.2.3` → published as a stable version.
- **Dev/pre-release**: `v1.2.3-dev.1` → marked as `prerelease` on GitHub automatically.

---

## 📦 CI/CD

| Workflow | Trigger | What it runs |
|----------------|----------------|-------------------------------------------|
| `ci.yml` | `push` or `PR` to `main` / `develop` | `make check` and matrix tests |
| `release.yml` | `push` a tag starting with `v` | `make check-release` + GitHub Release |

## 🫧 Before You PR

✅ Format, lint, test  
✅ Tidy modules  
✅ Never commit stray `replace` statements in adapters  
✅ Make sure CI passes

## 🙏 Thanks for contributing!

Open an issue or discussion if you have any questions.