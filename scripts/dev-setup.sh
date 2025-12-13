#!/bin/bash
set -e

echo "🚀 Starting development setup..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go 1.21 or later."
    exit 1
fi

echo "✅ Go version: $(go version)"

# Install dependencies
echo "📦 Installing dependencies..."
make deps

# Generate documentation
echo "📚 Generating documentation..."
make docs

# Run tests
echo "🧪 Running tests..."
make test

# Build the application
echo "🔨 Building application..."
make build

echo "✅ Development setup complete!"
echo ""
echo "📖 Available commands:"
echo "  make run              - Start development server"
echo "  make test             - Run all tests"
echo "  make test-coverage    - Generate coverage report"
echo "  make build            - Build the application"
echo "  make clean            - Clean build artifacts"
echo ""
echo "🌐 Once started, the API will be available at:"
echo "  http://localhost:8080/api/v1/users"
echo "  http://localhost:8080/swagger/"
echo "  http://localhost:8080/health"