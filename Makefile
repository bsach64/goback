BINARY=goback

all: build

tidy:
	go mod tidy

format: tidy
	gofmt -s -w .

lint: format
	golangci-lint run ./...

build: lint
	go build -v .

test: lint
	go test -v ./...

clean:
	go clean
	rm -r .data/*
