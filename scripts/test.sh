#!/bin/bash
set -e

echo "🧪 Running comprehensive test suite..."

# Clean test cache
echo "🧹 Cleaning test cache..."
go clean -testcache

# Run unit tests
echo "📝 Running unit tests..."
make test-unit

# Run integration tests
echo "🔗 Running integration tests..."
make test-integration

# Run race condition detection
echo "🏃‍♂️ Running race detection tests..."
make test-race

# Run benchmarks
echo "⚡ Running performance benchmarks..."
make benchmark

# Generate coverage report
echo "📊 Generating coverage report..."
make test-coverage

echo "✅ All tests completed successfully!"
echo ""
echo "📊 Coverage report: coverage/coverage.html"
echo "🌐 Open with: open coverage/coverage.html"