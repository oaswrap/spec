# Release Process

This document outlines the steps for creating and deploying new versions of `oaswrap/spec`. The entire release process is automated using GitHub Actions and is triggered by pushing a new Git tag.

## ## Prerequisites

Before creating a release, ensure your local environment is correctly set up:
1.  Check out the `main` branch: `git checkout main`
2.  Pull the latest changes from the remote: `git pull origin main`

## ## Development Release 🧪

Development releases are used for pre-releases or internal testing versions (e.g., `v1.2.3-dev.1`). They are marked as "pre-release" on GitHub.

### ### Step 1: Run Pre-Release Checks (Recommended)

Verify that the codebase passes all quality gates before tagging. This command runs tests, linting, and other essential checks.

```bash
make check-release
````

### ### Step 2: Create and Push the Dev Tag

Use the `release-dev` Makefile target. This single command runs final checks, creates the main tag and all associated adapter tags, and pushes them to `origin`.

Replace `vX.Y.Z-dev.N` with your actual version number.

```bash
make release-dev VERSION=v1.2.3-dev.1
```

Pushing the tag automatically triggers the **🚀 Release** workflow on GitHub, which builds and publishes the pre-release.

## ## Stable Release 🎉

Stable releases are official, production-ready versions (e.g., `v1.2.3`). They are marked as the "Latest" release on GitHub.

### ### Step 1: Finalize Dependencies and Commit

Ensure all internal dependencies are updated to the final release version. For example, if you were using dev versions in adapter modules, sync them to the new stable version number.

```bash
# Example: Update all adapter go.mod files
make sync-adapter-deps VERSION=v1.2.3

# Commit the final changes
git add .
git commit -m "build: prepare for release v1.2.3"
git push origin main
```

### ### Step 2: Run Final Release Checks

Run the comprehensive quality gate on the final commit to ensure it's stable.

```bash
make check-release
```

### ### Step 3: Create and Push the Stable Tag

Use the `release` Makefile target to kick off the production deployment.

Replace `vX.Y.Z` with your actual version number.

```bash
make release VERSION=v1.2.3
```

Pushing the stable tag triggers the **🚀 Release** workflow, which creates the official GitHub Release and marks it as **"Latest"**.

## ## Utility Commands

Your `Makefile` includes helpful commands for managing releases safely.

* **Dry Run**: To see what a release would do without creating or pushing any tags, use the `release-dry-run` target:

    ```bash
    make release-dry-run VERSION=v1.2.3
    ```

* **Fixing a Mistake**: To delete a mistakenly created version tag from both your local repository and the remote (`origin`), use the `delete-version` target:

    ```bash
    make delete-version VERSION=v1.2.3
    ```