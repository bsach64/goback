#!/bin/bash

set -e

echo "Formatting Go code..."
gofmt -s -w $(pwd)

echo "Running go mod tidy..."
go mod tidy

echo "Running tests..."
go test -v ./...

echo "------------Running golangcli-lint---------------"
echo "Running GoLint"
golangci-lint run
LINT_EXIT_CODE=$?

if [ $LINT_EXIT_CODE -ne 0 ]; then
    echo "Error During Linting. Fix the issues"
    exit $LINT_EXIT_CODE
fi

echo "------------All tests passed successfully---------------"

echo "Building the project"
go build ./...

echo "Tasks completed successfully!"

export PATH=$PATH:$(pwd) 
