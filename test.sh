#!/bin/bash

# Redis Repository Unit Test Runner
echo "🧪 Running Redis Repository Unit Tests"
echo "======================================"

echo ""
echo "📊 Running Unit Tests..."
go test ./internal/core/executor/ -v

echo ""
echo "✅ Unit tests completed!"
