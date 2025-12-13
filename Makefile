.PHONY: test test-unit test-integration test-coverage benchmark test-race deps clean lint docs run build

# Test targets
test: test-unit test-integration

test-unit:
	@echo "Running unit tests..."
	go test -v -short ./internal/store/... ./internal/handlers/...

test-integration:
	@echo "Running integration tests..."
	go test -v -run Integration ./internal/...

test-coverage:
	@echo "Running tests with coverage..."
	@mkdir -p coverage
	go test -v -coverprofile=coverage/coverage.out ./internal/...
	go tool cover -html=coverage/coverage.out -o coverage/coverage.html
	@echo "Coverage report generated: coverage/coverage.html"

test-race:
	@echo "Running tests with race detection..."
	go test -v -race ./internal/...

benchmark:
	@echo "Running benchmarks..."
	go test -v -bench=. -benchmem ./internal/store/...

# Dependencies
deps:
	go mod tidy
	go get github.com/stretchr/testify@v1.9.0
	go get gopkg.in/yaml.v3

# Linting
lint:
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run

# Documentation and build (existing targets enhanced)
docs:
	@which swag > /dev/null || (echo "Installing swag..." && go install github.com/swaggo/swag/cmd/swag@latest)
	swag init -g cmd/api-server/main.go

run: docs
	go run ./cmd/api-server

build: docs
	go build -o bin/api-server ./cmd/api-server

clean:
	rm -rf docs/ bin/ coverage/
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