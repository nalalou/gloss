#!/bin/bash
# Simulates a realistic AI agent fixing a bug
# Run: ./examples/agent-demo.sh | ./gloss watch

echo "::status id=read pending Reading files"
echo "::status id=fix pending Fix bug"
echo "::status id=test pending Run tests"
echo "::status id=commit pending Commit"
sleep 0.3

echo "::status id=read running Reading files"
echo "Reading src/auth/login.ts..."
sleep 0.4
echo "Reading src/auth/middleware.ts..."
sleep 0.4
echo "Reading src/auth/tokens.ts..."
sleep 0.3
echo "::status id=read done Read 3 files"
sleep 0.2

echo "::status id=fix running Analyzing"
echo ""
echo "Found the bug: token expiry comparison on line 42"
echo "compares timestamp as string instead of number."
sleep 0.8
echo ""
echo "Fixed: parseInt(expiry) > Date.now()"
sleep 0.4
echo "::status id=fix done Fix applied"
sleep 0.2

echo "::status id=test running Tests (0/3)"
echo ""
echo "Running auth/login tests..."
sleep 0.5
echo "  PASS TestLogin (0.3s)"
echo "  PASS TestLogout (0.1s)"
echo "::status id=test running Tests (1/3)"
sleep 0.5
echo "Running auth/middleware tests..."
echo "  PASS TestAuthMiddleware (0.2s)"
echo "::status id=test running Tests (2/3)"
sleep 0.5
echo "Running auth/tokens tests..."
echo "  PASS TestTokenRefresh (0.4s)"
echo "  PASS TestTokenExpiry (0.1s)"
echo "::status id=test done Tests (3/3, 5 passed)"
sleep 0.2

echo "::status id=commit running Committing"
sleep 0.4
echo "::status id=commit done Committed (abc1234)"

echo ""
echo "::ok Fix applied and verified — all 5 tests passing"
