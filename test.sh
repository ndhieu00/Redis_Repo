#!/bin/bash

# Redis Repository Unit Test Runner
echo "🧪 Running Redis Repository Unit Tests"
echo "======================================"

echo ""
echo "📊 Running Unit Tests..."
echo "Unit Test for command executor"
go test ./internal/core/executor/

echo ""
echo "✅ Unit tests completed!"
