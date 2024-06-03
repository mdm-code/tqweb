GO=go
GOFLAGS=-mod=vendor
COV_PROFILE=coverage.txt
MODULES=server,cmd

export CGO_ENABLED=0

.DEFAULT_GOAL := build

.PHONY: fmt vet lint test install build cover clean

fmt:
	@$(GO) fmt ./...

vet: fmt
	@$(GO) vet ./...

lint: vet
	@golint -set_exit_status=1 ./{$(MODULES)}/...

test: lint
	@$(GO) clean -testcache
	@$(GO) test ./... -v

install: test
	@$(GO) install ./...

build: test
	@$(GO) build $(GOFLAGS) github.com/mdm-code/tqweb/...

cover:
	@$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	@$(GO) tool cover -html=$(COV_PROFILE)

clean:
	@$(GO) clean github.com/mdm-code/tqweb/...
	@$(GO) mod tidy
	@$(GO) clean -testcache
	@rm -f $(COV_PROFILE)
