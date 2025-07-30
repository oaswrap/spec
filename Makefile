# Variables
PKG := ./...
COVERAGE_FILE := coverage.out

# Default target
.PHONY: all
all: test

# Run tests with gotestsum
.PHONY: test
test:
	@echo "Running tests with gotestsum..."
	@gotestsum --format standard-quiet -- $(PKG)

# Run tests and update golden files
.PHONY: test-update
test-update:
	@echo "Running tests with gotestsum and updating golden files..."
	@gotestsum --format standard-quiet -- -update $(PKG)

# Run tests with coverage and generate report
.PHONY: testcov
testcov:
	@echo "Running tests with gotestsum and coverage..."
	@gotestsum --format standard-quiet -- -coverprofile=$(COVERAGE_FILE) $(PKG)
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "Open HTML coverage report: make testcov-html"

# Open HTML coverage report
.PHONY: testcov-html
testcov-html:
	@go tool cover -html=$(COVERAGE_FILE)

# Format code
.PHONY: fmt
fmt:
	@echo "Formatting..."
	@go fmt $(PKG)

# Tidy go.mod and go.sum
.PHONY: tidy
tidy:
	@echo "Tidying modules..."
	@go mod tidy

# Lint code (requires golangci-lint)
.PHONY: lint
lint:
	@echo "Linting..."
	@golangci-lint run

# Clean up generated files
.PHONY: clean
clean:
	@echo "Cleaning..."
	@rm -f $(COVERAGE_FILE)

# Update dependencies
.PHONY: update
update:
	@echo "Updating dependencies..."
	@go get -u ./...