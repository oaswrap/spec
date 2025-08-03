# -------------------------------
# Vars
PKG := ./...
COVERAGE_FILE := coverage.out
ADAPTERS := fiberopenapi ginopenapi echoopenapi

# -------------------------------
# Core + Adapters: Tests
.PHONY: test
test:
	@echo "🔍 Core tests..."
	@gotestsum --format standard-quiet -- $(PKG)
	@for a in $(ADAPTERS); do \
		echo "🔍 Adapter $$a..."; \
		cd adapters/$$a && gotestsum --format standard-quiet -- ./... || exit 1; \
		cd ../..; \
	done

.PHONY: test-parallel
test-parallel:
	@gotestsum --format standard-quiet -- $(PKG) &
	@for a in $(ADAPTERS); do \
		(cd adapters/$$a && gotestsum --format standard-quiet -- ./...) & \
	done; \
	wait

# -------------------------------
# Coverage
.PHONY: testcov
testcov:
	@gotestsum --format standard-quiet -- -coverprofile=$(COVERAGE_FILE) $(PKG)
	@go tool cover -func=$(COVERAGE_FILE)
	@for a in $(ADAPTERS); do \
		cd adapters/$$a && gotestsum --format standard-quiet -- -coverprofile=$$a-$(COVERAGE_FILE) ./... && go tool cover -func=$$a-$(COVERAGE_FILE) || exit 1; \
		cd ../..; \
	done

# -------------------------------
# Tidy, Sync, Clean
.PHONY: tidy
tidy:
	@echo "🧹 Tidying core..."
	@go mod tidy
	@for a in $(ADAPTERS); do \
		echo "🧹 Tidying adapters/$$a..."; \
		cd adapters/$$a && go mod tidy && cd ../..; \
	done

.PHONY: sync
sync:
	@echo "🔗 Syncing workspace..."
	@go work sync

.PHONY: clean
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(COVERAGE_FILE)
	@for a in $(ADAPTERS); do rm -f adapters/$$a/$$a-$(COVERAGE_FILE); done

# -------------------------------
# Replace Management
.PHONY: fix-replace
fix-replace:
	@echo "🔧 Removing local replaces..."
	@for a in $(ADAPTERS); do \
		cd adapters/$$a && \
		if grep -q "replace github.com/oaswrap/spec" go.mod; then \
			go mod edit -dropreplace github.com/oaswrap/spec; \
			go mod tidy; \
			echo "✅ Removed replace in adapters/$$a"; \
		else \
			echo "✅ No replace in adapters/$$a"; \
		fi; \
		cd ../..; \
	done

.PHONY: check-replace-strict
check-replace-strict:
	@echo "🔍 Checking for accidental replaces..."
	@for a in $(ADAPTERS); do \
		if grep -q "replace github.com/oaswrap/spec" adapters/$$a/go.mod 2>/dev/null; then \
			echo "🚫 Found replace in adapters/$$a/go.mod"; \
			echo "💡 Run 'make fix-replace' to auto-fix"; \
			exit 1; \
		else \
			echo "✅ No replace in adapters/$$a/go.mod"; \
		fi; \
	done

# -------------------------------
# Lint & Tools
.PHONY: lint
lint:
	@echo "🔍 Linting core..."
	@golangci-lint run
	@for a in $(ADAPTERS); do \
		echo "🔍 Linting adapters/$$a..."; \
		cd adapters/$$a && golangci-lint run && cd ../..; \
	done

.PHONY: install-tools
install-tools:
	@echo "📦 Installing tools..."
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# -------------------------------
# Quality Gates
.PHONY: check
check: sync tidy lint test

.PHONY: check-release
check-release: sync tidy fix-replace lint test check-replace-strict

# -------------------------------
# Release Checks
.PHONY: release-check
release-check:
	@git diff --exit-code || (echo "❌ Uncommitted changes" && exit 1)
	@git diff --cached --exit-code || (echo "❌ Staged changes" && exit 1)
	@git status --porcelain | grep -q . && (echo "❌ Untracked files" && exit 1) || true

# -------------------------------
# Bump Dev
.PHONY: bump-dev
bump-dev:
ifndef NEXT
	$(error Usage: make bump-dev NEXT=v0.2.0-dev.1)
endif
	@echo "🔢 Bumping adapters to $(NEXT)..."
	@for a in $(ADAPTERS); do \
		cd adapters/$$a && \
		sed -i.bak -E 's#(github.com/oaswrap/spec )v[0-9]+\.[0-9]+\.[^ ]*#\1$(NEXT)#' go.mod && \
		rm -f go.mod.bak && \
		go mod tidy; \
		cd ../..; \
	done

# -------------------------------
# Release & Dev Tag
.PHONY: release
release: release-check
ifndef VERSION
	$(error Usage: make release VERSION=v1.2.3)
endif
	@echo "🚀 Running release quality gate..."
	@make check-release
	@echo "🏷️  Tagging release $(VERSION)..."
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "🎉 Production release $(VERSION) created and pushed!"

.PHONY: release-dev
release-dev: release-check
ifndef VERSION
	$(error Usage: make release-dev VERSION=v1.2.3-dev.1)
endif
	@echo "🚀 Running dev release checks..."
	@make check-release
	@echo "🏷️  Tagging dev release $(VERSION)..."
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "🎉 Dev release $(VERSION) created and pushed!"

# -------------------------------
# Help
.PHONY: help
help:
	@echo "make test                # Core + adapters tests"
	@echo "make test-parallel       # Parallel test"
	@echo "make testcov             # Coverage"
	@echo "make tidy                # go mod tidy"
	@echo "make sync                # go work sync"
	@echo "make clean               # Clean coverage"
	@echo "make lint                # Run linters"
	@echo "make fix-replace         # Drop local replaces"
	@echo "make check               # Local dev check"
	@echo "make check-release       # Full release check"
	@echo "make release-check       # Ensure clean git state"
	@echo "make bump-dev NEXT=...   # Bump adapters version"
	@echo "make release VERSION=... # Tag production release"
	@echo "make release-dev VERSION=...  # Tag dev version"