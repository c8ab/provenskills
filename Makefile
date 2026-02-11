.PHONY: build test lint clean

build:
	CGO_ENABLED=0 go build -ldflags="-s -w" -o psk ./cmd/psk/

test:
	go test ./...

lint:
	go vet ./...
	golangci-lint run

clean:
	rm -f psk
