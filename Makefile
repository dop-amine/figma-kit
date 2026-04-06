BINARY := figma-kit
PKG    := github.com/amine/figma-kit
CMD    := ./cmd/figma-kit
VERSION ?= dev
LDFLAGS := -ldflags "-s -w -X $(PKG)/internal/cli.Version=$(VERSION)"

.PHONY: build test lint install clean

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

test:
	go test -race ./...

lint:
	golangci-lint run ./...

install:
	go install $(LDFLAGS) $(CMD)

clean:
	rm -f $(BINARY)
