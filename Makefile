.PHONY: build install test clean

MOD = github.com/blazejsewera/go-test-proxy
BIN = gotestproxy
TO_BUILD = "$(MOD)/cmd/$(BIN)"

PREFIX ?= /usr/local/bin

all: build test

build:
	go build $(TO_BUILD)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build $(TO_BUILD)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build $(TO_BUILD)

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 go build $(TO_BUILD)

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build $(TO_BUILD)

install: all
	install -m755 "$(BIN)" "$(PREFIX)"

uninstall:
	rm -f "$(PREFIX)/$(BIN)"

test:
	go test ./...

clean:
	rm -f "$(BIN)"
	go clean -cache -testcache
