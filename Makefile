# ----------------------------------------
# oaswrap/spec Makefile
# ----------------------------------------

# Variables
PKG := ./...
COVERAGE_FILE := coverage.out
ADAPTERS := fiberopenapi ginopenapi echoopenapi

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
	@gotestsum --format standard-quiet -- $(go list ./... | grep -v '/adapters/')

.PHONY: test-adapters
test-adapters:
	@echo "🔍 Running tests for all adapters..."
	@for adapter in $(ADAPTERS); do \
		echo "🔍 Running tests for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- ./... || exit 1; \
		cd ../..; \
	done

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

.PHONY: test-update-core
test-update-core:
	@echo "📝 Updating golden files for core only..."
	@gotestsum --format standard-quiet -- -update $(go list ./... | grep -v '/adapters/')

.PHONY: test-update-adapters
test-update-adapters:
	@echo "📝 Updating golden files for all adapters..."
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
		cd adapters/$$adapter && gotestsum --format standard-quiet -- -coverprofile=$(COVERAGE_FILE) ./... && go tool cover -func=$(COVERAGE_FILE) || exit 1; \
		cd ../..; \
	done

.PHONY: testcov-core
testcov-core:
	@echo "📊 Running coverage for core only..."
	@gotestsum --format standard-quiet -- -coverprofile=core-$(COVERAGE_FILE) $(go list ./... | grep -v '/adapters/')
	@go tool cover -func=core-$(COVERAGE_FILE)

.PHONY: testcov-adapters
testcov-adapters:
	@echo "📊 Running coverage for all adapters..."
	@for adapter in $(ADAPTERS); do \
		echo "📊 Coverage for adapters/$$adapter..."; \
		cd adapters/$$adapter && gotestsum --format standard-quiet -- -coverprofile=$(COVERAGE_FILE) ./... && go tool cover -func=$(COVERAGE_FILE) || exit 1; \
		cd ../..; \
	done

.PHONY: testcov-html-core
testcov-html-core:
	@go tool cover -html=$(COVERAGE_FILE)

.PHONY: testcov-html-adapters
testcov-html-adapters:
	@for adapter in $(ADAPTERS); do \
		echo "📊 Opening HTML coverage for adapters/$$adapter..."; \
		cd adapters/$$adapter && go tool cover -html=$(COVERAGE_FILE) || exit 1; \
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

.PHONY: lint-core
lint-core:
	@golangci-lint run

.PHONY: lint-adapters
lint-adapters:
	@for adapter in $(ADAPTERS); do \
		echo "🔍 Linting adapters/$$adapter..."; \
		cd adapters/$$adapter && golangci-lint run || exit 1; \
		cd ../..; \
	done

# ----------------------------------------
# Other Maintenance
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
	@echo "🧹 Cleaning generated coverage files..."
	@rm -f $(COVERAGE_FILE) core-$(COVERAGE_FILE)
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

.PHONY: check
check: sync lint test
	@echo "✅ All checks passed!"

.PHONY: release-prep
release-prep:
	@echo "🚀 Preparing for release..."
	@$(MAKE) tidy
	@$(MAKE) sync
	@$(MAKE) lint
	@$(MAKE) test
	@echo "✅ Ready for release!"

.PHONY: help
help:
	@echo ""
	@echo "🚀 Available Make Targets"
	@echo "=========================="
	@echo ""
	@echo "📋 Testing:"
	@echo "  test                     Run tests for all modules"
	@echo "  test-core                Run tests for core module only"
	@echo "  test-adapters            Run tests for all adapters"
	@echo ""
	@echo "🔄 Golden Files:"
	@echo "  test-update              Update golden files for all modules"
	@echo "  test-update-core         Update golden files for core only"
	@echo "  test-update-adapters     Update golden files for all adapters"
	@echo ""
	@echo "📊 Coverage:"
	@echo "  testcov                  Run coverage for all modules"
	@echo "  testcov-core             Run coverage for core only"
	@echo "  testcov-adapters         Run coverage for all adapters"
	@echo "  testcov-html-core        Open HTML coverage for core"
	@echo "  testcov-html-adapters    Open HTML coverage for adapters"
	@echo ""
	@echo "🔍 Code Quality:"
	@echo "  fmt                      Format code"
	@echo "  lint                     Lint all modules"
	@echo "  lint-core                Lint core only"
	@echo "  lint-adapters            Lint adapters"
	@echo ""
	@echo "🔧 Maintenance:"
	@echo "  tidy                     Tidy modules"
	@echo "  sync                     Sync workspace"
	@echo "  clean                    Clean coverage files"
	@echo "  update                   Update dependencies"
	@echo ""
	@echo "✅ Workflows:"
	@echo "  check                    Run sync, lint, test"
	@echo "  release-prep             Prepare for release"
	@echo ""
	@echo "🛠️  Setup:"
	@echo "  install-tools            Install dev tools"
	@echo ""