.PHONY: build run test lint clean

build:
	go build -o bin/terminaltask ./cmd/terminaltask

run: build
	./bin/terminaltask

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
