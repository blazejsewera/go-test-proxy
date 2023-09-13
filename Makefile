.PHONY: build test clean

build:
	go build github.com/blazejsewera/go-test-proxy/cmd/gotestproxy

test:
	go test ./...

clean:
	rm -f gotestproxy
	go clean -cache -testcache
