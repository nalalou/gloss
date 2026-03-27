#!/bin/bash
# Simulates a deploy pipeline
# Run: ./examples/deploy-demo.sh | ./gloss watch

echo "::status id=lint pending Lint"
echo "::status id=types pending Type check"
echo "::status id=test pending Unit tests"
echo "::status id=build pending Build"
echo "::status id=push pending Push image"
echo "::status id=deploy pending Deploy"
sleep 0.3

echo "::status id=lint running Linting"
sleep 0.6
echo "  eslint: 0 errors, 0 warnings"
echo "::status id=lint done Lint"

echo "::status id=types running Type check"
sleep 0.5
echo "  tsc --noEmit: 0 errors"
echo "::status id=types done Type check"

echo "::status id=test running Unit tests"
sleep 0.4
echo "  auth: 12 passed"
sleep 0.3
echo "  api: 34 passed"
sleep 0.3
echo "  db: 8 passed"
echo "::status id=test done Unit tests (54 passed)"

echo "::status id=build running Building"
sleep 0.5
echo "  next build completed in 12s"
echo "  Bundle: 412KB / 500KB budget"
echo "::status id=build done Build (12s)"

echo "::status id=push running Pushing image"
sleep 0.6
echo "  app:v2.4.1 -> us-east-1.ecr"
echo "::status id=push done Image pushed"

echo "::status id=deploy running Rolling out"
sleep 0.4
echo "  Pod 1/3 healthy"
sleep 0.3
echo "  Pod 2/3 healthy"
sleep 0.3
echo "  Pod 3/3 healthy"
echo "::status id=deploy done Deployed"

echo ""
echo "::ok v2.4.1 live on production"
