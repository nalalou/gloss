#!/bin/bash
# Simulates agent output with :: protocol directives
# Run: ./examples/agent-demo.sh | ./gloss watch

echo "::status id=lint pending Lint"
echo "::status id=test pending Test"
echo "::status id=deploy pending Deploy"
echo "::bar id=prog 0 Progress"
sleep 0.5

echo "::status id=lint running Linting..."
echo "Checking formatting..."
echo "Checking imports..."
sleep 1
echo "::status id=lint done Lint (0 errors)"
echo "::bar id=prog 33 Progress"

echo "::status id=test running Testing..."
echo "PASS auth_test.go (12 tests)"
echo "PASS api_test.go (8 tests)"
sleep 1
echo "PASS db_test.go (6 tests)"
echo "::status id=test done Test (26 passed)"
echo "::bar id=prog 66 Progress"

echo "::status id=deploy running Deploying..."
echo "Pushing image..."
echo "Rolling out pods..."
sleep 1
echo "::status id=deploy done Deploy complete"
echo "::bar id=prog 100 Progress"
echo "::kv id=meta Pods=3/3 | Latency=42ms | Region=us-east-1"

echo "::ok All green"
