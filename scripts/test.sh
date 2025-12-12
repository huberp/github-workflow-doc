#!/bin/bash
# Test script for workflowdocgen

set -e

echo "Running Go tests..."
go test -v -race -coverprofile=coverage.out -covermode=atomic -coverpkg=./pkg/... ./pkg/...

echo ""
echo "Coverage report:"
if [ -f coverage.out ]; then
    go tool cover -func=coverage.out
    
    # Check coverage threshold (70%)
    COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
    THRESHOLD=70
    
    # Convert to integer by removing decimal point (e.g., 94.4 -> 944)
    COVERAGE_INT=$(echo "$COVERAGE" | awk '{printf "%.0f", $1 * 10}')
    THRESHOLD_INT=$((THRESHOLD * 10))
    
    echo ""
    echo "Total coverage: ${COVERAGE}%"
    echo "Required threshold: ${THRESHOLD}.0%"
    
    if [ "$COVERAGE_INT" -lt "$THRESHOLD_INT" ]; then
        echo "ERROR: Test coverage ${COVERAGE}% is below required threshold ${THRESHOLD}.0%"
        exit 1
    fi
    
    echo "✓ Coverage check passed!"
else
    echo "ERROR: No coverage data available"
    exit 1
fi

echo ""
echo "Running integration tests..."

# Create temporary test files
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Build the binary
go build -o "$TEMP_DIR/workflowdocgen" ./cmd/workflowdocgen

# Test 1: Help command
HELP_OUTPUT=$("$TEMP_DIR/workflowdocgen" --help 2>&1)
if [[ ! "$HELP_OUTPUT" =~ "output" ]]; then
    echo "Help command test failed!"
    exit 1
fi

# Test 2: Generate documentation from current workflows
mkdir -p "$TEMP_DIR/workflows"
cp -r .github/workflows/*.yml "$TEMP_DIR/workflows/" 2>/dev/null || true

if [ -n "$(ls -A "$TEMP_DIR/workflows" 2>/dev/null)" ]; then
    "$TEMP_DIR/workflowdocgen" --workflows-dir "$TEMP_DIR/workflows" --output "$TEMP_DIR/output.md"
    if [ ! -f "$TEMP_DIR/output.md" ]; then
        echo "Documentation generation test failed!"
        exit 1
    fi
    echo "Generated documentation:"
    head -20 "$TEMP_DIR/output.md"
fi

echo ""
echo "All integration tests passed! ✓"
