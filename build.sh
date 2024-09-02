#!/bin/bash

set -e

echo "Formatting Go code..."
gofmt -s -w $(pwd)

echo "Running go mod tidy..."
go mod tidy

echo "Running tests..."
go test -v ./...

echo "------------All tests passed successfully---------------"

echo "Building the project"
go build ./...

echo "Tasks completed successfully!"

