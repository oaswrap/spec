# Release Guide

This document outlines the release process for the oaswrap/spec project, which follows a multi-module architecture with a core module and multiple adapter modules.

Version input note: release commands accept both `x.y.z` and `vx.y.z` for `VERSION`.

Release model note: local development can use monorepo `replace` workflows, but released adapter tags must pin `github.com/oaswrap/spec` to the released core tag version.

## Project Structure

The project consists of:

- **Core module**: `github.com/oaswrap/spec` - The main OpenAPI specification builder
- **Adapter modules**: Framework-specific integrations
  - `github.com/oaswrap/spec/adapter/chiopenapi` - Chi framework adapter
  - `github.com/oaswrap/spec/adapter/echoopenapi` - Echo framework adapter
  - `github.com/oaswrap/spec/adapter/echov5openapi` - Echo v5 framework adapter
  - `github.com/oaswrap/spec/adapter/fiberopenapi` - Fiber framework adapter
  - `github.com/oaswrap/spec/adapter/ginopenapi` - Gin framework adapter
  - `github.com/oaswrap/spec/adapter/httpopenapi` - net/http adapter
  - `github.com/oaswrap/spec/adapter/httprouteropenapi` - HttpRouter adapter
  - `github.com/oaswrap/spec/adapter/muxopenapi` - Gorilla Mux adapter

## Prerequisites

Before releasing, ensure you have:

1. **Required tools installed**:
   ```bash
   make install-tools
   ```

2. **Clean working directory**:
   ```bash
   git status
   ```

3. **All tests passing**:
   ```bash
   make check
   ```

4. **Updated dependencies**:
   ```bash
   make tidy
   ```

## Release Types

### 1. Two-Stage Monorepo Release

Recommended end-to-end flow for root + adapters:

```bash
# Stage 1: release core and sync adapter deps
make release-prepare VERSION=x.y.z

# Commit sync changes produced by stage 1
git add adapter/*/go.mod adapter/*/go.sum
git commit -m "chore: sync adapter deps to vx.y.z"
git push

# Stage 2: publish adapter tags
make release-publish VERSION=x.y.z
```

Preview without changes:

```bash
make release-dry-run VERSION=x.y.z
```

## Typical Release Workflow

### Patch Release (x.y.Z)

For bug fixes and minor improvements:

1. **Prepare the release**:
   ```bash
   # Ensure everything is clean and tested
   make check
   
   # Update any necessary documentation
   # Update CHANGELOG.md if maintained
   ```

2. **Run Stage 1 (core tag + dependency sync)**:
   ```bash
   make release-prepare VERSION=0.3.5

   # Commit the sync result produced by stage 1
   git add adapter/*/go.mod adapter/*/go.sum
   git commit -m "chore: sync adapter deps to v0.3.5"
   git push
   
   # Publish adapter tags (stage 2)
   make release-publish VERSION=0.3.5
   ```

### Minor Release (x.Y.z)

For new features and non-breaking changes:

1. **Complete development and testing**:
   ```bash
   make check
   ```

2. **Update documentation** (README, examples, etc.)

3. **Run Stage 1 (core tag + dependency sync)**:
   ```bash
   make release-prepare VERSION=0.4.0

   # Commit the sync result produced by stage 1
   git add adapter/*/go.mod adapter/*/go.sum
   git commit -m "chore: sync adapter deps to v0.4.0"
   git push

   # Publish adapter tags (stage 2)
   make release-publish VERSION=0.4.0
   ```

### Major Release (X.y.z)

For breaking changes:

1. **Update migration guides** and documentation
2. **Thoroughly test** all modules and adapters
3. **Follow the same process** as minor releases but with careful version coordination

## Dependency Management

### Syncing Adapter Dependencies

When releasing a new core version, adapters need to be updated to reference the new version:

```bash
# Update all adapters to use the specified core version
make sync-adapter-deps VERSION=v1.2.0

# Skip go mod tidy during sync (useful for CI)
make sync-adapter-deps VERSION=v1.2.0 NO_TIDY=1

# Equivalent (also accepted):
make sync-adapter-deps VERSION=1.2.0
make sync-adapter-deps VERSION=1.2.0 NO_TIDY=1
```

### Cleaning Replace Directives

Remove local replace directives from adapter go.mod files:

```bash
make clean-replaces
```

## Version Management

### Semantic Versioning

The project follows [Semantic Versioning](https://semver.org/):

- **PATCH** (0.0.x): Bug fixes, documentation updates
- **MINOR** (0.x.0): New features, backwards-compatible changes
- **MAJOR** (x.0.0): Breaking changes

### Pre-release Testing (RC)

Use release candidate (RC) versions for testing before final release:

- `0.4.0-rc.1`
- `0.4.0-rc.2`

Example with two-stage release flow:

```bash
make release-dry-run VERSION=0.4.0-rc.1
make release-prepare VERSION=0.4.0-rc.1

# Commit sync changes from stage 1
git add adapter/*/go.mod adapter/*/go.sum
git commit -m "chore: sync adapter deps to v0.4.0-rc.1"
git push

make release-publish VERSION=0.4.0-rc.1
```

### Tag Management

#### Deleting Tags

If you need to delete a tag (use with caution):

```bash
make delete-tag TAG=v1.2.0
```

This will:
- Delete the local tag
- Delete the remote tag
- Require confirmation before deletion

#### Tag Naming Convention

- **Core module**: `v1.2.3`
- **Adapter modules**: `adapter/{name}/v1.2.3`

Examples:
- Core: `v0.3.5`
- Gin adapter: `adapter/ginopenapi/v0.3.5`
- Echo adapter: `adapter/echoopenapi/v0.3.5`

## CI/CD Integration

Current automation in this repository:

1. **Branch and pull-request pushes** trigger CI checks
2. **Automated testing** runs on all supported Go versions
3. **Tag-driven release publishing** is currently handled by Makefile commands and git tags

## Quality Checks

Before any release, ensure:

### 1. Code Quality
```bash
# Run linting
make lint

# Run all tests
make test

# Generate coverage reports
make testcov-html
```

### 2. Module Health
```bash
# Tidy all modules
make tidy-all

# Sync workspace
make sync

# List adapter status
make list-adapters
```

### 3. Documentation
- Update README.md if API changes
- Update examples if necessary
- Verify all adapter READMEs are current

## Troubleshooting

### Common Issues

1. **Tag already exists**:
   ```bash
   # Preview tag availability first
   make release-dry-run VERSION=1.2.0
   ```

2. **Adapter dependency mismatch**:
   ```bash
   # Resync dependencies
   make sync-adapter-deps VERSION=1.2.0
   ```

3. **Test failures**:
   ```bash
   # Update golden files if needed
   make test-update
   ```

4. **Module tidy issues**:
   ```bash
   # Clean and retidy all modules
   make tidy-all
   ```

## Release Checklist

### Pre-release
- [ ] All tests pass (`make test`)
- [ ] Linting passes (`make lint`)
- [ ] Documentation updated
- [ ] CHANGELOG.md updated (if maintained)
- [ ] Version number decided
- [ ] Clean git working directory

### Core Release
- [ ] `make release-prepare VERSION=x.y.z` completed successfully
- [ ] Tag appears in GitHub releases
- [ ] Module available on pkg.go.dev

### Adapter Release (if needed)
- [ ] Sync changes committed and pushed
- [ ] All adapter tests pass
- [ ] `make release-publish VERSION=x.y.z` completed successfully
- [ ] All adapter tags created

### Post-release
- [ ] Verify releases on GitHub
- [ ] Check pkg.go.dev for module availability
- [ ] Update any dependent projects
- [ ] Announce release (if applicable)

## Contact

For questions about the release process, please:
- Open an issue in the repository
- Check existing documentation
- Review the Makefile for available commands

## Additional Resources

- [Semantic Versioning](https://semver.org/)
- [Go Modules Reference](https://golang.org/ref/mod)
- [GitHub Releases Documentation](https://docs.github.com/en/repositories/releasing-projects-on-github)