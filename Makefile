BIN := copia
VERSION ?= DEV
DATE := $(shell date -u +%Y-%m-%d)
LDFLAGS := -s -w \
	-X github.com/qubernetic-org/copia-cli/internal/build.Version=$(VERSION) \
	-X github.com/qubernetic-org/copia-cli/internal/build.Date=$(DATE)

.PHONY: build test integration acceptance docs clean

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) ./cmd/copia

test:
	go test ./...

integration:
	go test -tags=integration ./...

acceptance:
	go test -tags=acceptance ./acceptance/...

docs:
	go run script/gen-docs.go

clean:
	rm -rf bin/ site/
