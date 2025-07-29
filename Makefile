# Variables
PKG := ./...
COVERAGE_FILE := coverage.out

# Default target
.PHONY: all
all: test

# Run tests
.PHONY: test
test:
	@echo "Running tests..."
	@go test -v $(PKG)

# Run tests with and update golden files
.PHONY: test-update
test-update:
	@echo "Running tests and updating golden files..."
	@go test -v $(PKG) -update

# Run tests with coverage and generate report
.PHONY: coverage
coverage:
	@echo "Running tests with coverage..."
	@go test -coverprofile=$(COVERAGE_FILE) $(PKG)
	@go tool cover -func=$(COVERAGE_FILE)
	@echo "Open HTML coverage report: make coverage-html"

# Open HTML coverage report
.PHONY: coverage-html
coverage-html:
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
