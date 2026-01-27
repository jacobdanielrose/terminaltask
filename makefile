.PHONY: build run test lint clean

BUILD_VERSION ?= $(shell git describe --tags --always --dirty="-dirty")
COMMIT ?= $(shell git rev-parse --short HEAD)
DATE   ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')


build:
	mkdir -p bin
	CGO_ENABLED=0 go build \
	-trimpath \
	-ldflags="-s -w \
	-X main.version=$(BUILD_VERSION) \
	-X main.commit=$(COMMIT) \
	-X main.buildDate=$(DATE)" \
	-o bin/terminaltask ./cmd/terminaltask

run: build
	./bin/terminaltask

release:
	@if [[ -z "$(VERSION)" ]]; then \
	    echo "Usage: make release VERSION=v0.1.1"; \
		exit 1; \
	fi
	git tag $(VERSION)
	git push origin $(VERSION)

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -rf bin/
