GO ?= /home/pi/.local/go/bin/go
BINARY := bin/coai

.PHONY: build test fmt

build:
	mkdir -p bin
	$(GO) build -o $(BINARY) ./cmd/coai

test:
	$(GO) test ./...

fmt:
	$(GO) fmt ./...
