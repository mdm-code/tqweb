GO=go
GOFLAGS=-mod=vendor -trimpath
COV_PROFILE=coverage.txt
TEMPL=templ

export CGO_ENABLED=0

.DEFAULT_GOAL := build

.PHONY: fmt vet test install build cover clean serve templ

templ:
	@$(TEMPL) generate

fmt: templ
	@$(GO) fmt ./...

vet: fmt
	@$(GO) vet ./...

test: vet
	@$(GO) clean -testcache
	@$(GO) test ./... -v

install: test
	@$(GO) install ./...

build: test
	@$(GO) build -C ./cmd/tqweb $(GOFLAGS) -o ../../tqweb

cover:
	@$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	@$(GO) tool cover -html=$(COV_PROFILE)

clean:
	@$(GO) clean github.com/mdm-code/tqweb/...
	@$(GO) mod tidy
	@$(GO) clean -testcache
	@rm -f $(COV_PROFILE)

serve: build
	@./tqweb
