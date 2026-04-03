BIN := copia-cli
VERSION ?= DEV
DATE := $(shell date -u +%Y-%m-%d)
LDFLAGS := -s -w \
	-X github.com/qubernetic/copia-cli/internal/build.Version=$(VERSION) \
	-X github.com/qubernetic/copia-cli/internal/build.Date=$(DATE)

.PHONY: build test integration acceptance docs docs-serve docs-clean clean snapshot install uninstall update

build:
	go build -ldflags "$(LDFLAGS)" -o bin/$(BIN) ./cmd/copia-cli

test:
	go test ./...

integration:
	go test -tags=integration ./...

acceptance:
	go test -tags=acceptance ./acceptance/...

docs:
	go run script/gen-docs.go

docs-serve: docs
	cd docs/site && bundle exec jekyll serve --baseurl /copia-cli

docs-clean:
	rm -rf docs/site/_site docs/site/.jekyll-cache docs/site/manual/copia-cli_*.md docs/site/_includes/sidebar.html

clean:
	rm -rf bin/ dist/

snapshot:
	goreleaser release --snapshot --clean

install: ## Install rpm (Fedora/RHEL only)
	sudo dnf install -y dist/$(BIN)_*_linux_amd64.rpm

uninstall:
	sudo dnf remove -y $(BIN)

update: snapshot
	sudo dnf remove -y $(BIN) || true
	sudo dnf install -y dist/$(BIN)_*_linux_amd64.rpm
