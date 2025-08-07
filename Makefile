# -------------------------------
# Vars
#
# Use := for immediate expansion, which is generally safer in Makefiles.
PKG           := ./...
COVERAGE_DIR  := coverage
COVERAGE_FILE := coverage.out
ADAPTERS      := chiopenapi echoopenapi fiberopenapi ginopenapi httpopenapi
MODULES	      := specui

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
.PHONY: test test-adapter test-update testcov testcov-html
.PHONY: tidy sync lint check
.PHONY: install-tools
.PHONY: list-adapters adapter-status
.PHONY: sync-adapter-deps
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
	@echo "$(BLUE)📋 Running additional tests...$(NC)"
	@for m in $(MODULES); do \
		echo "$(BLUE)🔍 Testing module $$m...$(NC)"; \
		(cd "module/$$m" && gotestsum --format standard-quiet -- ./...) || (echo "$(RED)❌ Module $$m tests failed$(NC)" && exit 1); \
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

	@for m in $(MODULES); do \
		echo "$(BLUE)📈 Module $$m coverage:$(NC)"; \
		(cd "module/$$m" && gotestsum --format standard-quiet -- -covermode=atomic -coverprofile="../../$(COVERAGE_DIR)/$$m-$(COVERAGE_FILE)" ./...); \
		if [ -f $(COVERAGE_DIR)/$$m-$(COVERAGE_FILE) ]; then \
			tail -n +2 $(COVERAGE_DIR)/$$m-$(COVERAGE_FILE) >> $(COVERAGE_DIR)/coverage.out; \
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
	@for m in $(MODULES); do \
		echo "$(BLUE)🧹 Tidying module/$$m...$(NC)"; \
		(cd "module/$$m" && go mod tidy); \
	done
	@echo "$(GREEN)✅ All modules tidied!$(NC)"

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
		(cd "adapter/$$a" && golangci-lint run) || (echo "$(RED)❌ Adapter $$a linting failed$(NC)" && exit 1); \
	done
	@echo "$(GREEN)🎉 All linting passed!$(NC)"

check: sync tidy lint test ## Run all local development checks
	@echo "$(GREEN)🎉 All local development checks passed!$(NC)"

install-tools: ## Install development tools
	@echo "$(BLUE)📦 Installing development tools...$(NC)"
	@go install gotest.tools/gotestsum@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)✅ Tools installed successfully!$(NC)"


list-adapters: ## List available adapters
	@echo "$(BLUE)📋 Available adapters:$(NC)"
	@for a in $(ADAPTERS); do echo "  - $$a"; done

adapter-status: ## Check the status of each adapter
	@echo "$(BLUE)📊 Adapter status overview:$(NC)"
	@for a in $(ADAPTERS); do \
		if [ -d "adapter/$$a" ]; then \
			echo "$(GREEN)✅ $$a$(NC) - exists"; \
		else \
			echo "$(RED)❌ $$a$(NC) - missing"; \
		fi; \
	done

sync-adapter-deps: ## Sync adapter dependencies to a specific version
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