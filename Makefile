.PHONY: build run clean install

build:
	go build -o polyterm .

run:
	go run main.go

install:
	go install

clean:
	rm -f polyterm

deps:
	go mod download
	go mod tidy

test:
	go test ./...

