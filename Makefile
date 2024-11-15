BINARY=goback

all: build

poner: 	#install but in spanish
 	go install

tidy:
	go mod tidy

format: tidy
	gofmt -s -w .

lint: format
	golangci-lint run ./...

build: lint
	go build -v .

database:
	sqlite3 -init createdb.sql meta.db .quit

test: lint
	go test -v ./...

clean:
	go clean
	rm -r .data/*
	rm meta.db
