# -------------------------------
# Vars
#
# Use := for immediate expansion, which is generally safer in Makefiles.
PKG           := ./...
COVERAGE_FILE := coverage.out
ADAPTERS      := chiopenapi echoopenapi fiberopenapi ginopenapi httpopenapi

# Platform detection for sed compatibility
# Using an immediately expanded variable for this is good practice.
UNAME_S := $(shell uname -s)
ifeq ($(UNAME_S),Darwin)
	SED_INPLACE := sed -i ''
else
	SED_INPLACE := sed -i
endif

# Colors for better output
RED    := \033[0;31m
GREEN  := \033[0;32m
YELLOW := \033[1;33m
BLUE   := \033[0;34m
NC     := \033[0m # No Color

# Ensure all targets are marked as phony to avoid conflicts with filenames.
.PHONY: test test-parallel test-update test-update-parallel
.PHONY: testcov coverage-html
.PHONY: tidy sync clean
.PHONY: fix-replace check-replace-strict
.PHONY: lint install-tools
.PHONY: list-adapters adapter-status
.PHONY: check check-release check-dry-run
.PHONY: release-check
.PHONY: release release-dry-run release-dry-run-clean release-status
.PHONY: check-adapter-deps sync-adapter-deps
.PHONY: list-tags delete-version verify-tags
.PHONY: help

# -------------------------------
# Core + Adapters: Tests
#
test:
	@echo "$(BLUE)🔍 Running core tests...$(NC)"
	@gotestsum --format standard-quiet -- $(PKG) || (echo "$(RED)❌ Core tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Core tests passed$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🔍 Testing adapter $$a...$(NC)"; \
		(cd "adapters/$$a" && gotestsum --format standard-quiet -- ./...) || (echo "$(RED)❌ Adapter $$a tests failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All tests passed!$(NC)"

# FIX: Added a final 'wait' and failure checking to ensure the script exits
# with an error if any parallel test fails.
test-parallel:
	@echo "$(BLUE)🚀 Running tests in parallel...$(NC)"
	@pids=""; \
	gotestsum --format standard-quiet -- $(PKG) & pids="$$pids $$!"; \
	for a in $(ADAPTERS); do \
		(cd "adapters/$$a" && gotestsum --format standard-quiet -- ./...) & pids="$$pids $$!"; \
	done; \
	\
	status=0; \
	for pid in $$pids; do \
		wait $$pid || status=1; \
	done; \
	if [ $$status -ne 0 ]; then \
		echo "$(RED)❌ One or more parallel tests failed$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)🎉 All parallel tests completed!$(NC)"

test-update:
	@echo "$(YELLOW)🔍 Running core tests (updating golden files)...$(NC)"
	@gotestsum --format standard-quiet -- -update $(PKG) || (echo "$(RED)❌ Core test update failed$(NC)" && exit 1)
	@for a in $(ADAPTERS); do \
		echo "$(YELLOW)🔍 Updating adapter $$a golden files...$(NC)"; \
		(cd "adapters/$$a" && gotestsum --format standard-quiet -- -update ./...) || (echo "$(RED)❌ Adapter $$a update failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)✅ All golden files updated!$(NC)"

# FIX: Added failure checking similar to 'test-parallel'.
test-update-parallel:
	@echo "$(YELLOW)🚀 Updating golden files in parallel...$(NC)"
	@pids=""; \
	gotestsum --format standard-quiet -- -update $(PKG) & pids="$$pids $$!"; \
	for a in $(ADAPTERS); do \
		(cd "adapters/$$a" && gotestsum --format standard-quiet -- -update ./...) & pids="$$pids $$!"; \
	done; \
	\
	status=0; \
	for pid in $$pids; do \
		wait $$pid || status=1; \
	done; \
	if [ $$status -ne 0 ]; then \
		echo "$(RED)❌ One or more parallel golden file updates failed$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ All golden files updated in parallel!$(NC)"

# -------------------------------
# Coverage
#
testcov:
	@echo "$(BLUE)📊 Generating coverage report...$(NC)"
	@gotestsum --format standard-quiet -- -coverprofile="coverage/$(COVERAGE_FILE)" $(PKG)
	@echo "$(BLUE)📈 Core coverage:$(NC)"
	@go tool cover -func="$(COVERAGE_FILE)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)📈 Adapter $$a coverage:$(NC)"; \
		(cd "adapters/$$a" && gotestsum --format standard-quiet -- -coverprofile="../../coverage/$$a-$(COVERAGE_FILE)" ./... && go tool cover -func="../../coverage/$$a-$(COVERAGE_FILE)") || exit 1; \
	done
	@echo "$(GREEN)✅ Coverage reports generated!$(NC)"

testcov-html: testcov
	@echo "$(BLUE)🌐 Generating HTML coverage reports...$(NC)"
	@go tool cover -html="coverage/$(COVERAGE_FILE)" -o "coverage/coverage.html"
	@for a in $(ADAPTERS); do \
		(cd "adapters/$$a" && go tool cover -html="../../coverage/$$a-$(COVERAGE_FILE)" -o "../../coverage/$$a-coverage.html"); \
	done
	@echo "$(GREEN)✅ HTML coverage reports generated!$(NC)"

# -------------------------------
# Tidy, Sync, Clean
#
tidy:
	@echo "$(BLUE)🧹 Tidying core...$(NC)"
	@go mod tidy
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🧹 Tidying adapters/$$a...$(NC)"; \
		(cd "adapters/$$a" && go mod tidy); \
	done
	@echo "$(GREEN)✅ All modules tidied!$(NC)"

sync:
	@echo "$(BLUE)🔗 Syncing workspace...$(NC)"
	@go work sync
	@echo "$(GREEN)✅ Workspace synced!$(NC)"

clean:
	@echo "$(BLUE)🧹 Cleaning coverage files...$(NC)"
	@rm -f $(COVERAGE_FILE) coverage.html
	@for a in $(ADAPTERS); do \
		rm -f "adapters/$$a/$$a-$(COVERAGE_FILE)" "adapters/$$a/$$a-coverage.html"; \
	done
	@echo "$(GREEN)✅ Cleanup completed!$(NC)"

# -------------------------------
# Replace Management
#
fix-replace:
	@echo "$(YELLOW)🔧 Removing local replaces...$(NC)"
	@for a in $(ADAPTERS); do \
		if grep -q "replace github.com/oaswrap/spec" "adapters/$$a/go.mod"; then \
			(cd "adapters/$$a" && go mod edit -dropreplace github.com/oaswrap/spec && go mod tidy); \
			echo "$(GREEN)✅ Removed replace in adapters/$$a$(NC)"; \
		else \
			echo "$(GREEN)✅ No replace needed in adapters/$$a$(NC)"; \
		fi; \
	done

check-replace-strict:
	@echo "$(BLUE)🔍 Checking for accidental replaces...$(NC)"
	@for a in $(ADAPTERS); do \
		if grep -q "replace github.com/oaswrap/spec" "adapters/$$a/go.mod" 2>/dev/null; then \
			echo "$(RED)🚫 Found replace in adapters/$$a/go.mod$(NC)"; \
			echo "$(YELLOW)💡 Run 'make fix-replace' to auto-fix$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(GREEN)✅ No accidental replaces found.$(NC)"

# -------------------------------
# Lint & Tools
#
lint:
	@echo "$(BLUE)🔍 Linting core...$(NC)"
	@golangci-lint run || (echo "$(RED)❌ Core linting failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Core linting passed$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🔍 Linting adapters/$$a...$(NC)"; \
		(cd "adapters/$$a" && golangci-lint run) || (echo "$(RED)❌ Adapter $$a linting failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All linting passed!$(NC)"

install-tools:
	@echo "$(BLUE)📦 Installing development tools...$(NC)"
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✅ Tools installed successfully!$(NC)"

# -------------------------------
# Utility
#
list-adapters:
	@echo "$(BLUE)📋 Available adapters:$(NC)"
	@for a in $(ADAPTERS); do echo "  - $$a"; done

adapter-status:
	@echo "$(BLUE)📊 Adapter status overview:$(NC)"
	@for a in $(ADAPTERS); do \
		if [ -d "adapters/$$a" ]; then \
			echo "$(GREEN)✅ $$a$(NC) - exists"; \
		else \
			echo "$(RED)❌ $$a$(NC) - missing"; \
		fi; \
	done

# -------------------------------
# Quality Gates
#
check: sync tidy lint test
	@echo "$(GREEN)🎉 All local development checks passed!$(NC)"

check-release: sync tidy check-replace-strict lint test
	@echo "$(GREEN)🎉 All release checks passed!$(NC)"

check-dry-run:
	@echo "$(YELLOW)🔍 Dry run - would execute: sync tidy lint test$(NC)"
	@echo "$(BLUE)Current git status:$(NC)"
	@git status --short

# -------------------------------
# Release Checks
#
release-check:
	@echo "$(BLUE)🔍 Checking git state for release...$(NC)"
	@if ! git diff --exit-code --quiet; then \
		echo "$(RED)❌ Uncommitted changes detected$(NC)"; \
		exit 1; \
	fi
	@if ! git diff --cached --exit-code --quiet; then \
		echo "$(RED)❌ Staged changes detected$(NC)"; \
		exit 1; \
	fi
	# FIX: Redirect stderr to /dev/null to silence "Not a git repository" errors in certain CI environments.
	@BRANCH=$$(\
		git branch --show-current 2>/dev/null \
	); \
	if [ "$$BRANCH" != "main" ] && [ "$$BRANCH" != "master" ]; then \
		echo "$(YELLOW)⚠️  Warning: Not on main/master branch (current: $$BRANCH)$(NC)"; \
	fi
	@echo "$(GREEN)✅ Git state is clean for release$(NC)"

# -------------------------------
# Release Management
#
release: release-check
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release VERSION=v0.3.0$(NC)"; \
		exit 1; \
	fi

	@echo "$(BLUE)🚀 Running release quality gate...$(NC)"
	@$(MAKE) check-release

	@echo "$(BLUE)🔄 Syncing adapter dependencies to $(VERSION)...$(NC)"
	@$(MAKE) sync-adapter-deps VERSION=$(VERSION) NO_TIDY=1

	@echo "$(BLUE)📥 Committing updated adapter dependencies...$(NC)"
	@git add .
	@git diff --cached --quiet || git commit -m "chore: sync adapters to $(VERSION)"

	@echo "$(BLUE)🏷️  Tagging main release $(VERSION)...$(NC)"
	@git tag -f "$(VERSION)"

	@echo "$(BLUE)🏷️  Tagging adapter releases...$(NC)"
	@if [ -n "$(ADAPTERS)" ]; then \
		for a in $(ADAPTERS); do \
			ADAPTER_TAG="adapters/$$a/$(VERSION)"; \
			git tag -f "$$ADAPTER_TAG"; \
			echo "$(GREEN)✅ Tagged $$ADAPTER_TAG$(NC)"; \
		done \
	fi

	@echo "$(BLUE)📤 Pushing commit and main tag first...$(NC)"
	@git push origin HEAD
	@git push origin "$(VERSION)"

	echo "$(BLUE)⏳ Waiting for Go proxy to index $(VERSION)...$(NC)"
	@sleep 5

	@echo "$(BLUE)🔍 Forcing proxy refresh for root module...$(NC)"
	@GOPROXY=proxy.golang.org go list -m github.com/oaswrap/spec@$(VERSION) || true

	@echo "$(BLUE)🔍 Fallback: direct fetch to ensure fresh tag is visible...$(NC)"
	@GOPROXY=direct go list -m github.com/your/repo@$(VERSION) || true

	@echo "$(BLUE)📤 Pushing adapter tags...$(NC)"
	@if [ -n "$(ADAPTERS)" ]; then \
		for a in $(ADAPTERS); do \
			ADAPTER_TAG="adapters/$$a/$(VERSION)"; \
			git push origin "$$ADAPTER_TAG"; \
			echo "$(GREEN)✅ Pushed $$ADAPTER_TAG$(NC)"; \
		done \
	fi

	@echo "$(BLUE)🧹 Tidying all modules now that tags are pushed...$(NC)"
	@$(MAKE) tidy
	@git add .
	@git diff --cached --quiet || git commit -m "chore: tidy modules after $(VERSION)"

	@echo "$(GREEN)✅ Tidy completed and committed after release push!$(NC)"
	@echo "$(GREEN)🎉 Production release $(VERSION) created and pushed!$(NC)"

# -------------------------------
# Development Release Management
#
release-dry-run: release-check
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-dry-run VERSION=v0.3.0$(NC)"; \
		exit 1; \
	fi

	@echo "$(BLUE)🚀 [Dry Run] Running release quality gate...$(NC)"
	@$(MAKE) check-release

	@echo "$(BLUE)🔄 [Dry Run] Syncing adapter dependencies to $(VERSION)...$(NC)"
	@$(MAKE) sync-adapter-deps VERSION=$(VERSION) NO_TIDY=1

	@echo "$(BLUE)📥 [Dry Run] Staged changes that would be committed:$(NC)"
	@git add .
	@git diff --cached --name-status

	@echo "$(BLUE)🏷️  [Dry Run] Tags that would be created:$(NC)"
	@echo "  - $(VERSION)"
	@if [ -n "$(ADAPTERS)" ]; then \
		for a in $(ADAPTERS); do \
			echo "  - adapters/$$a/$(VERSION)"; \
		done \
	fi

	@echo "$(BLUE)📤 [Dry Run] Would push in this order:$(NC)"
	@echo "  - Push HEAD"
	@echo "  - Push main tag: $(VERSION)"
	@echo "  - Wait & warm Go proxy"
	@echo "  - Force direct fetch fallback"
	@echo "  - Push adapter tags:"

	@if [ -n "$(ADAPTERS)" ]; then \
		for a in $(ADAPTERS); do \
			echo "    - adapters/$$a/$(VERSION)"; \
		done \
	fi

	@echo "$(BLUE)🧹 [Dry Run] Would tidy after pushing tags$(NC)"
	@echo "$(GREEN)✅ [Dry Run] Release plan looks good!$(NC)"

release-dry-run-clean:
	@git reset
	@git checkout -- .
	@echo "$(GREEN)✅ Cleaned staged and working tree changes from dry run.$(NC)"

release-status:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-status VERSION=v0.3.0$(NC)"; \
		exit 1; \
	fi

	@echo "$(BLUE)🔍 Checking available versions at proxy.golang.org...$(NC)"
	@echo "$(BLUE)Available versions for github.com/oaswrap/spec:$(NC)"
	@curl -s https://proxy.golang.org/github.com/oaswrap/spec/@v/list || echo "$(RED)❌ Could not reach proxy.golang.org$(NC)"

	@echo ""
	@echo "$(BLUE)🔍 Checking if $(VERSION) is indexed at proxy.golang.org...$(NC)"
	@curl -s https://proxy.golang.org/github.com/oaswrap/spec/@v/$(VERSION).info || echo "$(RED)❌ Version not found at proxy.golang.org$(NC)"

	@echo ""
	@echo "$(BLUE)🔍 Checking checksum database sum.golang.org for $(VERSION)...$(NC)"
	@curl -s https://sum.golang.org/lookup/github.com/oaswrap/spec@$(VERSION) || echo "$(RED)❌ Not found in checksum DB$(NC)"

	@echo ""
	@echo "$(BLUE)🔍 Testing go list -m at proxy.golang.org...$(NC)"
	@GOPROXY=proxy.golang.org go list -m github.com/oaswrap/spec@$(VERSION) || echo "$(RED)❌ Not found via go list -m (proxy)$(NC)"

	@echo ""
	@echo "$(BLUE)🔍 Testing go list -m with direct fallback...$(NC)"
	@GOPROXY=direct go list -m github.com/oaswrap/spec@$(VERSION) || echo "$(RED)❌ Not found via go list -m (direct)$(NC)"

	@echo ""
	@echo "$(GREEN)✅ Status check done for $(VERSION)!$(NC)"

# -------------------------------
# Dependency Management
#
# FIX: Escaped shell variables and improved logic for clarity.
check-adapter-deps:
	@echo "$(BLUE)🔍 Checking adapter dependencies...$(NC)"
	@LATEST_TAG=$$(\
		git tag -l 'v*' --sort=-version:refname | grep -v 'adapters/' | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$$' | head -1\
	); \
	for a in $(ADAPTERS); do \
		echo "$(BLUE)📋 Adapter: $$a$(NC)"; \
		SPEC_VERSION=$$(\
			grep 'github.com/oaswrap/spec ' "adapters/$$a/go.mod" | awk '{print $$2}' \
		); \
		echo "  References: $$SPEC_VERSION"; \
	done; \
	echo ""; \
	echo "$(BLUE)📋 Latest stable release tag: $$LATEST_TAG$(NC)"; \
	echo ""; \
	echo "$(YELLOW)💡 To sync all adapters to the latest release, run:$(NC)"; \
	echo "   make sync-adapter-deps VERSION=$$LATEST_TAG"

# -----------------------------------
# Sync adapter dependencies to an existing version
#
# Use this AFTER you publish a version.
# Example:
#   make sync-adapter-deps VERSION=v0.3.0
#
sync-adapter-deps:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make sync-adapter-deps VERSION=v0.3.0 [NO_TIDY=1]$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)🔄 Syncing adapter dependencies to $(VERSION)...$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)📝 Updating adapters/$$a...$(NC)"; \
		(cd "adapters/$$a" && \
		$(SED_INPLACE) -E 's#(github.com/oaswrap/spec )v[0-9]+\.[0-9]+\.[^ ]*#\1$(VERSION)#' go.mod); \
		if [ "$(NO_TIDY)" != "1" ]; then \
			(cd "adapters/$$a" && go mod tidy); \
		else \
			echo "$(YELLOW)⚠️  Skipped go mod tidy for adapters/$$a because NO_TIDY=1$(NC)"; \
		fi; \
		echo "$(GREEN)✅ Updated adapters/$$a to $(VERSION)$(NC)"; \
	done
	@echo "$(GREEN)🎉 All adapters synced to $(VERSION)!$(NC)"

# -------------------------------
# Tag Management
#
list-tags:
	@echo "$(BLUE)📋 All version tags:$(NC)"
	@git tag -l 'v*' --sort=-version:refname | head -10
	@echo ""
	@echo "$(BLUE)📋 Adapter tags for latest version:$(NC)"
	@LATEST=$$(\
		git tag -l 'v*' --sort=-version:refname | grep -v 'adapters/' | head -1\
	); \
	if [ -n "$$LATEST" ]; then \
		echo "Latest version: $$LATEST"; \
		for a in $(ADAPTERS); do \
			if git tag -l "adapters/$$a/$$LATEST" | grep -q .; then \
				echo "$(GREEN)✅ adapters/$$a/$$LATEST$(NC)"; \
			else \
				echo "$(RED)❌ adapters/$$a/$$LATEST$(NC)"; \
			fi; \
		done; \
	else \
		echo "No version tags found"; \
	fi

delete-version:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make delete-version VERSION=v1.2.3$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)⚠️  This will delete version $(VERSION) and all related adapter tags!$(NC)"
	@printf "Press Enter to continue or Ctrl+C to cancel... "; \
	read -r
	@echo "$(BLUE)🗑️  Deleting local tags...$(NC)"
	@git tag -d "$(VERSION)" 2>/dev/null || true
	@for a in $(ADAPTERS); do \
		git tag -d "adapters/$$a/$(VERSION)" 2>/dev/null || true; \
	done
	@echo "$(BLUE)🗑️  Deleting remote tags...$(NC)"
	@git push --delete origin "$(VERSION)" 2>/dev/null || true
	@for a in $(ADAPTERS); do \
		git push --delete origin "adapters/$$a/$(VERSION)" 2>/dev/null || true; \
	done
	@echo "$(GREEN)✅ Version $(VERSION) deleted locally and remotely!$(NC)"

verify-tags:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make verify-tags VERSION=v1.2.3$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)🔍 Verifying tags for version $(VERSION)...$(NC)"
	@if git tag -l "$(VERSION)" | grep -q .; then \
		echo "$(GREEN)✅ Main tag $(VERSION) exists$(NC)"; \
	else \
		echo "$(RED)❌ Main tag $(VERSION) missing$(NC)"; \
	fi
	@for a in $(ADAPTERS); do \
		ADAPTER_TAG="adapters/$$a/$(VERSION)"; \
		if git tag -l "$$ADAPTER_TAG" | grep -q .; then \
			echo "$(GREEN)✅ $$ADAPTER_TAG exists$(NC)"; \
		else \
			echo "$(RED)❌ $$ADAPTER_TAG missing$(NC)"; \
		fi; \
	done

# -------------------------------
# Help
#
help:
	@echo "$(BLUE)Available targets:$(NC)"
	@echo ""
	@echo "$(YELLOW)Testing & Coverage:$(NC)"
	@echo "  test                     Run all tests sequentially"
	@echo "  test-parallel            Run all tests in parallel"
	@echo "  test-update              Update golden test files"
	@echo "  testcov                  Generate coverage reports"
	@echo "  coverage-html            Generate HTML coverage reports"
	@echo ""
	@echo "$(YELLOW)Development Workflow:$(NC)"
	@echo "  install-tools            Install required dev tools (gotestsum, golangci-lint)"
	@echo "  lint                     Run linters on all modules"
	@echo "  tidy                     Run 'go mod tidy' on all modules"
	@echo "  sync                     Run 'go work sync'"
	@echo "  clean                    Clean up generated files"
	@echo ""
	@echo "$(YELLOW)Quality Gates & Pre-release Checks:$(NC)"
	@echo "  check                    Run local dev checks (sync, tidy, lint, test)"
	@echo "  check-release            Run stricter checks for a release"
	@echo "  release-check            Ensure git state is clean for a release"
	@echo "  fix-replace              Remove local 'replace' directives from go.mod files"
	@echo "  check-replace-strict     Fail if any local 'replace' directives are found"
	@echo ""
	@echo "$(YELLOW)Release & Version Management:$(NC)"
	@echo "  release VERSION=...      Create and push a new production release tag"
	@echo "  release-dry-run VERSION=...  Dry run of the release process"
	@echo "  delete-version VERSION=..Delete a version tag locally and remotely"
	@echo ""
	@echo "$(YELLOW)Utilities & Information:$(NC)"
	@echo "  help                     Show this help message"
	@echo "  list-adapters            List all configured adapters"
	@echo "  list-tags                List recent tags and check adapter tag coverage"
	@echo "  verify-tags VERSION=...  Check if all tags for a version exist"
	@echo "  check-adapter-deps       Check the main dependency version for each adapter"
	@echo "  sync-adapter-deps VERSION=..Sync all adapters to a specific main dependency version"
