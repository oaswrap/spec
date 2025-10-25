# Release Guide

This document outlines the release process for the oaswrap/spec project, which follows a multi-module architecture with a core module and multiple adapter modules.

## Project Structure

The project consists of:

- **Core module**: `github.com/oaswrap/spec` - The main OpenAPI specification builder
- **Adapter modules**: Framework-specific integrations
  - `github.com/oaswrap/spec/adapter/chiopenapi` - Chi framework adapter
  - `github.com/oaswrap/spec/adapter/echoopenapi` - Echo framework adapter
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

### 1. Core Module Release

Release the main `github.com/oaswrap/spec` module:

```bash
make release VERSION=x.y.z
```

**Example:**
```bash
make release VERSION=1.2.0
```

This will:
- Create and push a Git tag `vx.y.z`
- Trigger the release workflow

### 2. Adapter Modules Release

Release all adapter modules simultaneously:

```bash
make release-adapters VERSION=x.y.z
```

**Example:**
```bash
make release-adapters VERSION=1.2.0
```

This will:
- Create tags for each adapter: `adapter/{name}/vx.y.z`
- Push all adapter tags to the repository

### 3. Dry Run for Adapters

Test the adapter release process without making changes:

```bash
make release-adapters-dry-run VERSION=x.y.z
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

2. **Release core module**:
   ```bash
   make release VERSION=0.3.5
   ```

3. **Update adapters** (if needed):
   ```bash
   # Sync adapter dependencies to new core version
   make sync-adapter-deps VERSION=v0.3.5
   
   # Release adapters
   make release-adapters VERSION=0.3.5
   ```

### Minor Release (x.Y.z)

For new features and non-breaking changes:

1. **Complete development and testing**:
   ```bash
   make check
   ```

2. **Update documentation** (README, examples, etc.)

3. **Release core module**:
   ```bash
   make release VERSION=0.4.0
   ```

4. **Update and release adapters**:
   ```bash
   make sync-adapter-deps VERSION=v0.4.0
   make release-adapters VERSION=0.4.0
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

The release process integrates with GitHub Actions:

1. **Pushing tags** triggers release workflows
2. **Automated testing** runs on all supported Go versions
3. **Release artifacts** are generated automatically
4. **Go modules** are published to the Go module proxy

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
   # Delete the tag and try again
   make delete-tag TAG=v1.2.0
   make release VERSION=1.2.0
   ```

2. **Adapter dependency mismatch**:
   ```bash
   # Resync dependencies
   make sync-adapter-deps VERSION=v1.2.0
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
- [ ] `make release VERSION=x.y.z` completed successfully
- [ ] Tag appears in GitHub releases
- [ ] Module available on pkg.go.dev

### Adapter Release (if needed)
- [ ] Dependencies synced (`make sync-adapter-deps VERSION=vx.y.z`)
- [ ] All adapter tests pass
- [ ] `make release-adapters VERSION=x.y.z` completed successfully
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