#!/bin/bash
# Runs gloss's own test suite with :: protocol output
# Usage: ./examples/self-test.sh | ./gloss watch

GLOSS_DIR="$(cd "$(dirname "$0")/.." && pwd)"
cd "$GLOSS_DIR"

# Declare packages to test
packages=(
  "internal/env"
  "internal/font"
  "internal/render"
  "internal/theme"
  "internal/protocol"
  "internal/watch"
  "cmd"
)

total=${#packages[@]}
passed=0
failed=0

# Set up panel
for pkg in "${packages[@]}"; do
  short="${pkg##*/}"
  echo "::status id=$short pending $short"
done
echo "::bar id=prog 0 Testing"

# Run each package
for i in "${!packages[@]}"; do
  pkg="${packages[$i]}"
  short="${pkg##*/}"

  echo "::status id=$short running $short"

  # Run tests, capture output
  output=$(go test "./$pkg/..." -v 2>&1)
  exit_code=$?

  # Print test output to scroll zone
  echo "$output" | grep -E "^(=== RUN|--- PASS|--- FAIL|PASS|FAIL|ok)" | while read line; do
    echo "    $line"
  done

  if [ $exit_code -eq 0 ]; then
    count=$(echo "$output" | grep -c "^--- PASS")
    echo "::status id=$short done $short ($count tests)"
    passed=$((passed + count))
  else
    count=$(echo "$output" | grep -c "^--- FAIL")
    echo "::status id=$short error $short ($count failures)"
    failed=$((failed + count))
  fi

  pct=$(( (i + 1) * 100 / total ))
  echo "::bar id=prog $pct Testing"
done

# Summary
total_tests=$((passed + failed))
echo "::kv id=summary Packages=$total | Tests=$total_tests | Passed=$passed | Failed=$failed"

if [ $failed -eq 0 ]; then
  echo "::ok All $total_tests tests passing"
else
  echo "::err $failed tests failed"
fi
