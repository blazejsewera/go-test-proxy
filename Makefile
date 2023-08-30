PROXY=github.com/blazejsewera/go-test-proxy/cmd/gotestproxy

.PHONY: build test clean

build:
	go build $(PROXY)

test:
	go test ./...

clean:
	rm -f $(PROXY_EXE)
	go clean -cache -testcache
