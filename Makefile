# Variables
PKG := ./...
COVERAGE_FILE := coverage.out
ADAPTERS := fiberopenapi ginopenapi echoopenapi

# ----------------------------------------
# Default target
.PHONY: all
all: sync test

# ----------------------------------------
# Core + Adapters: Tests
# ----------------------------------------

.PHONY: test
test:
	@echo "🔍 Running tests for core..."
	@gotestsum --format standard-quiet -- $(PKG)
	@for adapter in $(ADAPTERS); do \
		echo "🔍 Running tests for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- ./... || exit 1; \
		cd ../..; \
	done

.PHONY: test-core
test-core:
	@echo "🔍 Running tests for core only..."
	@gotestsum --format standard-quiet -- `go list ./... | grep -v '/adapters/'`

.PHONY: test-adapters
test-adapters:
	@for adapter in $(ADAPTERS); do \
		echo "🔍 Running tests for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- ./... || exit 1; \
		cd ../..; \
	done

.PHONY: test-parallel
test-parallel:
	@echo "🚀 Running tests in parallel..."
	@gotestsum --format standard-quiet -- $(PKG) &
	@for adapter in $(ADAPTERS); do \
		(cd adapters/$$adapter && gotestsum --format standard-quiet -- ./...) & \
	done; \
	wait

# ----------------------------------------
# Core + Adapters: Tests + Golden update
# ----------------------------------------

.PHONY: test-update
test-update:
	@echo "📝 Running tests + updating golden files..."
	@gotestsum --format standard-quiet -- -update $(PKG)
	@for adapter in $(ADAPTERS); do \
		echo "📝 Updating golden files for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- -update ./... || exit 1; \
		cd ../..; \
	done

# ----------------------------------------
# Coverage
# ----------------------------------------

.PHONY: testcov
testcov:
	@echo "📊 Running tests with coverage..."
	@gotestsum --format standard-quiet -- -coverprofile=$(COVERAGE_FILE) $(PKG)
	@go tool cover -func=$(COVERAGE_FILE)
	@for adapter in $(ADAPTERS); do \
		echo "📊 Running coverage for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- -coverprofile=$$adapter-$(COVERAGE_FILE) ./... && go tool cover -func=$$adapter-$(COVERAGE_FILE) || exit 1; \
		cd ../..; \
	done

# ----------------------------------------
# Lint & Format
# ----------------------------------------

.PHONY: fmt
fmt:
	@echo "✨ Formatting core..."
	@go fmt $(PKG)
	@for adapter in $(ADAPTERS); do \
		echo "✨ Formatting adapters/$$adapter..."; \
		cd adapters/$$adapter && go fmt ./... || exit 1; \
		cd ../..; \
	done

.PHONY: lint
lint:
	@echo "🔍 Linting core..."
	@golangci-lint run
	@for adapter in $(ADAPTERS); do \
		echo "🔍 Linting adapters/$$adapter..."; \
		cd adapters/$$adapter && golangci-lint run || exit 1; \
		cd ../..; \
	done

# ----------------------------------------
# Maintenance
# ----------------------------------------

.PHONY: tidy
tidy:
	@echo "🧹 Tidying core..."
	@go mod tidy
	@for adapter in $(ADAPTERS); do \
		echo "🧹 Tidying adapters/$$adapter..."; \
		cd adapters/$$adapter && go mod tidy || exit 1; \
		cd ../..; \
	done

.PHONY: sync
sync:
	@echo "🔗 Syncing workspace..."
	@go work sync

.PHONY: clean
clean:
	@echo "🧹 Cleaning coverage files..."
	@rm -f $(COVERAGE_FILE)
	@for adapter in $(ADAPTERS); do \
		rm -f adapters/$$adapter/$$adapter-$(COVERAGE_FILE); \
	done

.PHONY: update
update:
	@echo "⬆️  Updating dependencies..."
	@go get -u ./...
	@for adapter in $(ADAPTERS); do \
		echo "⬆️  Updating adapters/$$adapter..."; \
		cd adapters/$$adapter && go get -u ./... || exit 1; \
		cd ../..; \
	done

.PHONY: install-tools
install-tools:
	@echo "📦 Installing dev tools..."
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

.PHONY: check-tools
check-tools:
	@echo "🔧 Checking required tools..."
	@command -v gotestsum >/dev/null 2>&1 || (echo "❌ gotestsum not found. Run 'make install-tools'" && exit 1)
	@command -v golangci-lint >/dev/null 2>&1 || (echo "❌ golangci-lint not found. Run 'make install-tools'" && exit 1)
	@echo "✅ All tools available"

# ----------------------------------------
# 🔧 Replace Statement Management
# ----------------------------------------

.PHONY: fix-replace
fix-replace:
	@echo "🔧 Removing replace statements from adapters..."
	@for adapter in $(ADAPTERS); do \
		echo "🔧 Checking adapters/$$adapter..."; \
		cd adapters/$$adapter && \
		if grep -q "replace github.com/oaswrap/spec" go.mod 2>/dev/null; then \
			echo "  ➜ Removing replace statement"; \
			go mod edit -dropreplace github.com/oaswrap/spec; \
			go mod tidy; \
		else \
			echo "  ✅ No replace statements found"; \
		fi; \
		cd ../..; \
	done
	@echo "✅ Replace cleanup complete!"

.PHONY: check-replace
check-replace:
	@echo "🔍 Checking for accidental 'replace' in adapters..."
	@if grep -R "replace github.com/oaswrap/spec" adapters/*/go.mod 2>/dev/null; then \
		echo "🚫 ERROR: Found 'replace' in adapter go.mod! Please remove or move to go.work only."; \
		echo "💡 Run 'make fix-replace' to auto-fix"; \
		exit 1; \
	else \
		echo "✅ No accidental replaces found."; \
	fi

.PHONY: check-replace-strict
check-replace-strict:
	@echo "🔍 Strict replace check with auto-fix option..."
	@for adapter in $(ADAPTERS); do \
		if grep -q "replace github.com/oaswrap/spec" adapters/$$adapter/go.mod 2>/dev/null; then \
			echo "🚫 Found replace in adapters/$$adapter/go.mod"; \
			echo "💡 Run 'make fix-replace' to auto-fix"; \
			exit 1; \
		fi; \
	done
	@echo "✅ No replace statements found"

# ----------------------------------------
# ✅ Quality Gates (Development - No Replace Cleanup)
# ----------------------------------------

.PHONY: check
check: check-tools sync tidy lint test
	@echo "✅ All development checks passed!"

.PHONY: check-dev
check-dev: check-tools sync tidy lint test
	@echo "✅ Development quality gate passed!"

# ----------------------------------------
# 🚀 Release Management
# ----------------------------------------

.PHONY: release-check
release-check:
	@echo "🔍 Pre-release checks..."
	@git diff --exit-code || (echo "❌ Uncommitted changes found" && exit 1)
	@git diff --cached --exit-code || (echo "❌ Staged changes found" && exit 1)
	@git status --porcelain | grep -q . && (echo "❌ Untracked files found" && exit 1) || true
	@echo "✅ Repository is clean"

.PHONY: release
release: release-check
ifndef VERSION
	$(error Usage: make release VERSION=v1.2.3)
endif
	@echo "🚀 Creating release $(VERSION)..."
	@echo "Running full quality checks with replace cleanup..."
	@make check-release
	@echo "✅ All checks passed!"
	@echo "🏷️  Creating and pushing tag..."
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "🎉 Release $(VERSION) created! Check GitHub Actions for progress."
	@echo "🔗 View at: https://github.com/$(shell git config --get remote.origin.url | sed 's/.*github.com[:/]\([^.]*\).*/\1/')/releases"

.PHONY: release-clean
release-clean: release-check
ifndef VERSION
	$(error Usage: make release-clean VERSION=v1.2.3)
endif
	@echo "🚀 Creating clean release $(VERSION)..."
	
	# Auto-fix any replace statements
	@make fix-replace
	
	# Check if any changes were made
	@if ! git diff --exit-code --quiet; then \
		echo "🔧 Auto-fixed replace statements"; \
		git add .; \
		git commit -m "chore: remove replace statements before release"; \
	fi
	
	# Continue with normal release
	@make check-release
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "🎉 Clean release $(VERSION) created!"

.PHONY: release-dev
release-dev: release-check
ifndef VERSION
	$(error Usage: make release-dev VERSION=v1.2.3-dev.1)
endif
	@echo "🧪 Creating dev release $(VERSION)..."
	@make test
	@make fix-replace
	@git tag $(VERSION)
	@git push origin $(VERSION)
	@echo "🎉 Dev release $(VERSION) created!"

# ----------------------------------------
# ✅ Release Quality Gate (With Replace Cleanup)
# ----------------------------------------

.PHONY: check-release
check-release: check-tools sync tidy fix-replace lint test check-replace-strict
	@echo "✅ All release checks passed with replace cleanup!"

# ----------------------------------------
# Dev Version Bumping
# ----------------------------------------

.PHONY: bump-dev
bump-dev:
ifndef NEXT
	$(error Usage: make bump-dev NEXT=v0.2.0-dev.1)
endif
	@echo "🔢 Bumping adapters to use core version $(NEXT)..."
	@for adapter in $(ADAPTERS); do \
		echo "🔢 Updating adapters/$$adapter to $(NEXT)..."; \
		cd adapters/$$adapter && \
		sed -i.bak -E 's#(github.com/oaswrap/spec )v[0-9]+\.[0-9]+\.[0-9]+.*#\1$(NEXT)#' go.mod && \
		rm -f go.mod.bak && \
		go mod tidy; \
		cd ../..; \
	done
	@echo "✅ Bump done. Check and commit!"

# ----------------------------------------
# 📚 Help
# ----------------------------------------

.PHONY: help
help:
	@echo ""
	@echo "🚀 Available Make Targets"
	@echo "=========================="
	@echo ""
	@echo "📋 Testing:"
	@echo "  test                     Run tests for all modules"
	@echo "  test-core                Run tests for core only"
	@echo "  test-adapters            Run tests for all adapters"
	@echo "  test-parallel            Run tests in parallel"  
	@echo "  test-update              Run tests + update golden files"
	@echo ""
	@echo "📊 Coverage:"
	@echo "  testcov                  Run coverage for all modules"
	@echo ""
	@echo "🔍 Code Quality:"
	@echo "  fmt                      Format code"
	@echo "  lint                     Lint all modules"
	@echo ""
	@echo "🔧 Maintenance:"
	@echo "  tidy                     Tidy all modules"
	@echo "  sync                     Sync workspace"
	@echo "  clean                    Clean generated coverage"
	@echo "  update                   Update dependencies"
	@echo "  install-tools            Install dev tools"
	@echo "  check-tools              Check if required tools are installed"
	@echo ""
	@echo "🔧 Replace Statement Management:"
	@echo "  fix-replace              Auto-remove replace statements from adapters"
	@echo "  check-replace            Check for accidental replace statements"
	@echo "  check-replace-strict     Strict replace check with fix suggestions"
	@echo ""
	@echo "✅ Quality Gates:"
	@echo "  check                    Development quality gate (no replace cleanup)"
	@echo "  check-dev                Alias for 'check' (development)"
	@echo "  check-release            Release quality gate (with replace cleanup)"
	@echo "  bump-dev NEXT=...        Bump adapters to next dev version"
	@echo ""
	@echo "🚀 Release Management:"
	@echo "  release VERSION=...      Create production release (with replace cleanup)"
	@echo "  release-clean VERSION=... Create release with auto replace cleanup"
	@echo "  release-dev VERSION=...  Create development release (pre-release)"
	@echo "  release-check            Check if repository is ready for release"
	@echo ""
	@echo "🔄 Development vs Release:"
	@echo "  • Development: 'make check' - Fast, no replace cleanup"
	@echo "  • Release: 'make check-release' - Full validation with replace cleanup"
	@echo "  • CI uses 'check-release' for release validation"
	@echo ""
	@echo "Examples:"
	@echo "  make check                                     # Fast development checks"
	@echo "  make release VERSION=v1.2.3                    # Production release"
	@echo "  make release-clean VERSION=v1.2.3              # Safe release with cleanup"
	@echo "  make release-dev VERSION=v1.2.3-dev.1          # Development release"
	@echo "  make fix-replace                               # Remove replace statements"
	@echo "  make bump-dev NEXT=v1.3.0-dev.1                # Bump adapters version"
	@echo ""