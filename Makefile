GO ?= /home/pi/.local/go/bin/go
BINARY := bin/moai-codex

.PHONY: build test fmt

build:
	mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/moai-codex

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...
