.PHONY: test test-unit test-integration test-coverage benchmark test-race deps clean lint docs run build

# Test targets
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	go test -v -short ./store/... ./handlers/...

test-integration:
	@echo "Running integration tests..."
	go test -v -run Integration ./...

test-coverage:
	@echo "Running tests with coverage..."
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./...

benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./store/...

# Dependencies
deps:
	go mod tidy
	go get github.com/stretchr/testify@v1.9.0

# Linting
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Documentation and build (existing targets enhanced)
docs:
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init

run: docs
	go run main.go

build: docs
	go build -o api-server main.go

clean:
	rm -rf docs/ api-server
	go clean -testcache
	rm -f coverage.out coverage.html

# Run all checks
ci: deps test test-race test-coverage lint
	@echo "All checks passed!"

# Help
help:
	@echo "Available targets:"
	@echo "  test           - Run all tests"
	@echo "  test-unit      - Run unit tests only"  
	@echo "  test-integration - Run integration tests only"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race detection"
	@echo "  benchmark      - Run performance benchmarks"
	@echo "  deps           - Install/update dependencies"
	@echo "  lint           - Run code linting"
	@echo "  docs           - Generate Swagger documentation"
	@echo "  run            - Run development server"
	@echo "  build          - Build the application"
	@echo "  clean          - Clean build artifacts and test cache"
	@echo "  ci             - Run all CI checks"
	@echo "  help           - Show this help message"