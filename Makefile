# -------------------------------
# Vars
#
# Use := for immediate expansion, which is generally safer in Makefiles.
PKG           := ./...
COVERAGE_DIR  := coverage
COVERAGE_FILE := coverage.out
ADAPTERS      := chiopenapi echoopenapi fiberopenapi ginopenapi httpopenapi muxopenapi httprouteropenapi echov5openapi fiberv3openapi

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

# Tool versions
GOLANGCI_LINT_VERSION := v2.3.1
GOTESTSUM_VERSION     := v1.12.3

# Normalize VERSION input so targets accept both 1.2.3 and v1.2.3.
# Pre-release versions are also supported, e.g. 0.4.0-rc.1 or v0.4.0-rc.1.
VERSION_STRIPPED := $(patsubst v%,%,$(VERSION))
VERSION_TAG      := v$(VERSION_STRIPPED)

# Ensure all targets are marked as phony to avoid conflicts with filenames.
.PHONY: test test-adapter test-update testcov testcov-html
.PHONY: tidy sync lint check tidy-all
.PHONY: install-tools
.PHONY: list-adapters adapter-status
.PHONY: sync-adapter-deps
.PHONY: release-preflight release-core-publish
.PHONY: release-adapters-preflight release-adapters-publish release-adapters-publish-dry-run
.PHONY: release-prepare release-publish release-dry-run
.PHONY: help

help: ## Show this help message
	@echo "$(BLUE)Available targets:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-20s$(NC) %s\n", $$1, $$2}'

test: ## Run all tests (core + adapters)
	@echo "$(BLUE)🔍 Running core tests...$(NC)"
	@gotestsum --format standard-quiet -- $(PKG) || (echo "$(RED)❌ Core tests failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Core tests passed$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🔍 Testing adapter $$a...$(NC)"; \
		(cd "adapter/$$a" && gotestsum --format standard-quiet -- ./...) || (echo "$(RED)❌ Adapter $$a tests failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All tests passed!$(NC)"

test-adapter: ## Run tests for all adapters
	@echo "$(BLUE)🔍 Running tests for all adapters...$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🔍 Testing adapter $$a...$(NC)"; \
		(cd "adapter/$$a" && gotestsum --format standard-quiet -- ./...) || (echo "$(RED)❌ Adapter $$a tests failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All adapter tests passed!$(NC)"

test-update: ## Update golden files for tests
	@echo "$(YELLOW)🔍 Running core tests (updating golden files)...$(NC)"
	@gotestsum --format standard-quiet -- -update $(PKG) || (echo "$(RED)❌ Core test update failed$(NC)" && exit 1)
	@for a in $(ADAPTERS); do \
		echo "$(YELLOW)🔍 Updating adapter $$a golden files...$(NC)"; \
		(cd "adapter/$$a" && gotestsum --format standard-quiet -- -update ./...) || (echo "$(RED)❌ Adapter $$a update failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)✅ All golden files updated!$(NC)"

testcov: ## Run tests with coverage and generate reports
	@echo "$(BLUE)📊 Generating coverage report...$(NC)"
	@mkdir -p $(COVERAGE_DIR)
	@gotestsum --format standard-quiet -- -covermode=atomic -coverprofile="$(COVERAGE_DIR)/$(COVERAGE_FILE)" $(PKG)

	@for a in $(ADAPTERS); do \
		echo "$(BLUE)📈 Adapter $$a coverage:$(NC)"; \
		(cd "adapter/$$a" && gotestsum --format standard-quiet -- -covermode=atomic -coverprofile="../../$(COVERAGE_DIR)/$$a-$(COVERAGE_FILE)" ./...); \
		if [ -f $(COVERAGE_DIR)/$$a-$(COVERAGE_FILE) ]; then \
			tail -n +2 $(COVERAGE_DIR)/$$a-$(COVERAGE_FILE) >> $(COVERAGE_DIR)/coverage.out; \
		fi; \
	done

	@echo "$(BLUE)📊 Combined coverage report saved to $(COVERAGE_DIR)/$(COVERAGE_FILE)$(NC)"
	@go tool cover -func="$(COVERAGE_DIR)/$(COVERAGE_FILE)"

testcov-html: testcov ## Generate HTML coverage reports
	@echo "$(BLUE)🌐 Generating HTML coverage reports...$(NC)"
	@go tool cover -html="coverage/$(COVERAGE_FILE)" -o "coverage/coverage.html"
	@echo "$(GREEN)✅ HTML coverage reports generated!$(NC)"
	@open coverage/coverage.html

tidy: ## Tidy up Go modules for core and adapters
	@echo "$(BLUE)🧹 Tidying core...$(NC)"
	@go mod tidy
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🧹 Tidying adapter/$$a...$(NC)"; \
		(cd "adapter/$$a" && go mod tidy); \
	done
	@echo "$(GREEN)✅ All modules tidied!$(NC)"

tidy-all: ## Tidy up Go modules for all submodules
	@echo "🧹 Running go mod tidy in all modules..."
	@find . -name "go.mod" -execdir sh -c 'echo "📦 Tidying $$(pwd)" && go mod tidy' \;
	@echo "✅ All modules tidied."

sync: ## Sync Go workspace
	@echo "$(BLUE)🔗 Syncing workspace...$(NC)"
	@go work sync
	@echo "$(GREEN)✅ Workspace synced!$(NC)"

lint: ## Run linting
	@echo "$(BLUE)🔍 Linting core...$(NC)"
	@golangci-lint run || (echo "$(RED)❌ Core linting failed$(NC)" && exit 1)
	@echo "$(GREEN)✅ Core linting passed$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)🔍 Linting adapter/$$a...$(NC)"; \
		golangci-lint run ./adapter/$$a/... || \
			(echo "$(RED)❌ Adapter $$a linting failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All linting passed!$(NC)"

check: sync tidy lint test ## Run all local development checks
	@echo "$(GREEN)🎉 All local development checks passed!$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)📦 Installing development tools...$(NC)"
	@go install gotest.tools/gotestsum@$(GOTESTSUM_VERSION)
	@go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION)
	@echo "$(GREEN)✅ Tools installed successfully!$(NC)"

list-adapters: ## List available adapters
	@echo "$(BLUE)📋 Available adapters:$(NC)"
	@for a in $(ADAPTERS); do \
		if [ -d "adapter/$$a" ]; then \
			echo "$(GREEN)✅ $$a$(NC) - exists"; \
		else \
			echo "$(RED)❌ $$a$(NC) - missing"; \
		fi; \
	done

release-core-publish: ## Internal: release core module tag
	@$(MAKE) release-preflight VERSION=$(VERSION_TAG)
	@echo "$(BLUE)🚀 Releasing version $(VERSION_TAG)...$(NC)"
	@git tag -a $(VERSION_TAG) -m "Release $(VERSION_TAG)"
	@git push origin $(VERSION_TAG)

release-preflight: ## Validate release prerequisites for core tag
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-prepare VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)❌ Working tree is not clean. Commit or stash changes first.$(NC)"; \
		exit 1; \
	fi
	@if ! git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then \
		echo "$(RED)❌ No upstream branch configured. Push your branch and set upstream first.$(NC)"; \
		exit 1; \
	fi
	@if [ "$$(git rev-list --count @{u}..HEAD)" -ne 0 ]; then \
		echo "$(RED)❌ You have unpushed commits. Push them before running release-prepare.$(NC)"; \
		exit 1; \
	fi
	@if git rev-parse -q --verify "refs/tags/$(VERSION_TAG)" >/dev/null; then \
		echo "$(RED)❌ Local tag $(VERSION_TAG) already exists.$(NC)"; \
		exit 1; \
	fi
	@if git ls-remote --exit-code --tags origin "refs/tags/$(VERSION_TAG)" >/dev/null 2>&1; then \
		echo "$(RED)❌ Remote tag $(VERSION_TAG) already exists on origin.$(NC)"; \
		exit 1; \
	fi
	@echo "$(GREEN)✅ Core release preflight passed for $(VERSION_TAG)$(NC)"

release-adapters-publish: ## Internal: release adapter module tags
	@$(MAKE) release-adapters-preflight VERSION=$(VERSION_TAG)
	@echo "$(BLUE)🚀 Releasing adapters with version $(VERSION_TAG)...$(NC)"
	@tags=""; \
	for a in $(ADAPTERS); do \
		tag="adapter/$$a/$(VERSION_TAG)"; \
		echo "$(BLUE)🏷️  Creating local tag $$tag...$(NC)"; \
		git tag -a "$$tag" -m "Release $$tag"; \
		tags="$$tags $$tag"; \
	done; \
	echo "$(BLUE)🚀 Pushing adapter tags to origin...$(NC)"; \
	git push origin $$tags
	@echo "$(GREEN)🎉 All adapters released with version $(VERSION_TAG)!$(NC)"

release-adapters-preflight: ## Validate release prerequisites for adapter tags
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-publish VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "$(RED)❌ Working tree is not clean. Commit or stash changes first.$(NC)"; \
		exit 1; \
	fi
	@if ! git rev-parse --abbrev-ref --symbolic-full-name @{u} >/dev/null 2>&1; then \
		echo "$(RED)❌ No upstream branch configured. Push your branch and set upstream first.$(NC)"; \
		exit 1; \
	fi
	@if [ "$$(git rev-list --count @{u}..HEAD)" -ne 0 ]; then \
		echo "$(RED)❌ You have unpushed commits. Push them before running release-publish.$(NC)"; \
		exit 1; \
	fi
	@for a in $(ADAPTERS); do \
		tag="adapter/$$a/$(VERSION_TAG)"; \
		if ! grep -Eq 'github.com/oaswrap/spec[[:space:]]+$(VERSION_TAG)($$|[[:space:]])' "adapter/$$a/go.mod"; then \
			echo "$(RED)❌ adapter/$$a/go.mod is not pinned to $(VERSION_TAG). Run make sync-adapter-deps VERSION=$(VERSION_TAG).$(NC)"; \
			exit 1; \
		fi; \
		if git rev-parse -q --verify "refs/tags/$$tag" >/dev/null; then \
			echo "$(RED)❌ Local tag $$tag already exists.$(NC)"; \
			exit 1; \
		fi; \
		if git ls-remote --exit-code --tags origin "refs/tags/$$tag" >/dev/null 2>&1; then \
			echo "$(RED)❌ Remote tag $$tag already exists on origin.$(NC)"; \
			exit 1; \
		fi; \
	done
	@echo "$(GREEN)✅ Adapter release preflight passed for $(VERSION_TAG)$(NC)"

release-adapters-publish-dry-run:
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-dry-run VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)🔍 Dry run for releasing adapters with version $(VERSION_TAG)...$(NC)"
	@for a in $(ADAPTERS); do \
		tag="adapter/$$a/$(VERSION_TAG)"; \
		status="ok"; \
		if git rev-parse -q --verify "refs/tags/$$tag" >/dev/null; then status="local-tag-exists"; fi; \
		if git ls-remote --exit-code --tags origin "refs/tags/$$tag" >/dev/null 2>&1; then status="remote-tag-exists"; fi; \
		echo "$(BLUE)🚀 Would release adapter $$a with version $$tag (status: $$status)$(NC)"; \
	done
	@echo "$(GREEN)🎉 Dry run complete! No changes made.$(NC)"

release-prepare: ## Stage 1 release flow: root release + adapter dependency sync
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-prepare VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)📦 Stage 1/2: Preparing monorepo release for $(VERSION_TAG)...$(NC)"
	@$(MAKE) release-core-publish VERSION=$(VERSION_TAG)
	@$(MAKE) sync-adapter-deps VERSION=$(VERSION_TAG)
	@echo "$(YELLOW)⚠️  Commit and push adapter dependency changes before publishing adapter tags.$(NC)"
	@echo "$(YELLOW)   Suggested commit: chore: sync adapter deps to $(VERSION_TAG)$(NC)"
	@echo "$(GREEN)✅ Stage 1 complete. Next: make release-publish VERSION=$(VERSION_TAG)$(NC)"

release-publish: ## Stage 2 release flow: publish adapter tags from committed sync state
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-publish VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)📦 Stage 2/2: Publishing adapter tags for $(VERSION_TAG)...$(NC)"
	@$(MAKE) release-adapters-publish VERSION=$(VERSION_TAG)
	@echo "$(GREEN)🎉 Monorepo release completed for $(VERSION_TAG)!$(NC)"

release-dry-run: ## Dry run for two-stage monorepo release
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make release-dry-run VERSION=0.3.0 (or v0.3.0)$(NC)"; \
		exit 1; \
	fi
	@echo "$(YELLOW)🔍 Dry run for monorepo release $(VERSION_TAG)...$(NC)"
	@echo "$(BLUE)Stage 1/2 would create and push core tag: $(VERSION_TAG)$(NC)"
	@if git rev-parse -q --verify "refs/tags/$(VERSION_TAG)" >/dev/null; then \
		echo "$(YELLOW)   Local status: tag already exists$(NC)"; \
	else \
		echo "$(GREEN)   Local status: tag available$(NC)"; \
	fi
	@if git ls-remote --exit-code --tags origin "refs/tags/$(VERSION_TAG)" >/dev/null 2>&1; then \
		echo "$(YELLOW)   Remote status: tag already exists$(NC)"; \
	else \
		echo "$(GREEN)   Remote status: tag available$(NC)"; \
	fi
	@echo "$(BLUE)Stage 1/2 would sync adapter go.mod dependencies to $(VERSION_TAG)$(NC)"
	@$(MAKE) release-adapters-publish-dry-run VERSION=$(VERSION_TAG)
	@echo "$(YELLOW)ℹ️  Release process note: commit synced adapter dependencies before Stage 2 publishing.$(NC)"

delete-tag: ## Delete a Git tag
ifndef TAG
	$(error TAG is undefined. Usage: make delete-tag TAG=tagname)
endif
	@echo "Are you sure you want to delete tag $(TAG)? [y/N]" && read ans && [ $${ans:-N} = y ]
	git tag -d $(TAG)
	git push origin :refs/tags/$(TAG)
	@echo "Tag $(TAG) deleted successfully"

sync-adapter-deps: ## Sync adapter dependencies
	@if [ -z "$(VERSION)" ]; then \
		echo "$(RED)Usage: make sync-adapter-deps VERSION=0.3.0 (or v0.3.0) [NO_TIDY=1]$(NC)"; \
		exit 1; \
	fi
	@echo "$(BLUE)🔄 Syncing adapter dependencies to $(VERSION_TAG)...$(NC)"
	@for a in $(ADAPTERS); do \
		echo "$(BLUE)📝 Updating adapter/$$a...$(NC)"; \
		(cd "adapter/$$a" && \
		$(SED_INPLACE) -E 's#(github.com/oaswrap/spec )v[0-9]+\.[0-9]+\.[^ ]*#\1$(VERSION_TAG)#' go.mod); \
		if [ "$(NO_TIDY)" != "1" ]; then \
			(cd "adapter/$$a" && go mod tidy); \
		else \
			echo "$(YELLOW)⚠️  Skipped go mod tidy for adapter/$$a because NO_TIDY=1$(NC)"; \
		fi; \
		echo "$(GREEN)✅ Updated adapter/$$a to $(VERSION_TAG)$(NC)"; \
	done
	@echo "$(GREEN)🎉 All adapters synced to $(VERSION_TAG)!$(NC)"

.PHONY: clean-replaces
clean-replaces: ## Clean up replace directives in go.mod adapters
	@find adapter -mindepth 2 -maxdepth 2 -type f -name go.mod \
		-exec sed -i.bak '/^replace github\.com\/oaswrap\/spec =>/d' {} \; \
		-exec rm {}.bak \; \
		-execdir go mod tidy \;