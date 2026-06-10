.PHONY: build install test clean

MOD = github.com/blazejsewera/go-test-proxy
BIN = gotestproxy

PREFIX ?= /usr/local/bin

all: build test

build:
	go build -o $(BIN) $(MOD)

build-windows-amd64:
	GOOS=windows GOARCH=amd64 go build -o $(BIN).exe $(MOD)

build-linux-amd64:
	GOOS=linux GOARCH=amd64 go build -o $(BIN) $(MOD)

build-macos-amd64:
	GOOS=darwin GOARCH=amd64 go build -o $(BIN) $(MOD)

build-macos-arm64:
	GOOS=darwin GOARCH=arm64 go build -o $(BIN) $(MOD)

install: all
	install -m755 "$(BIN)" "$(PREFIX)"

uninstall:
	rm -f "$(PREFIX)/$(BIN)"

test:
	go test ./...

clean:
	rm -f "$(BIN)"
	go clean -cache -testcache
