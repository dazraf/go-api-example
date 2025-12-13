#!/bin/bash
set -e

echo "🔨 Building Go API Server..."

# Clean previous builds
echo "🧹 Cleaning previous builds..."
make clean

# Install dependencies
echo "📦 Installing dependencies..."
make deps

# Generate documentation
echo "📚 Generating API documentation..."
make docs

# Run tests
echo "🧪 Running tests..."
make test

# Build binary
echo "🏗️ Building binary..."
make build

echo "✅ Build complete!"
echo "📦 Binary location: bin/api-server"
echo ""
echo "🚀 To run:"
echo "  ./bin/api-server"
echo "  or"
echo "  make run"