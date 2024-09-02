BINARY=goback

all: build

tidy:
	go mod tidy

format: tidy
	gofmt -s -w .

lint: format
	golangci-lint run ./... || (echo "Lint errors generated..")

build: lint
	go build -v .

test: lint
	go test -v ./...

clean:
	go clean
	rm tmp/*
