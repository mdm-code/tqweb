GO=go
GOFLAGS=-mod=vendor
COV_PROFILE=coverage.txt

export CGO_ENABLED=0

.DEFAULT_GOAL := build

.PHONY: fmt vet test install build cover clean serve update-js gen-templ

fmt: gen-templ
	@$(GO) fmt ./...

vet: fmt
	@$(GO) vet ./...

test: vet
	@$(GO) clean -testcache
	@$(GO) test ./... -v

install: test
	@$(GO) install ./...

build: test update-js
	@$(GO) build $(GOFLAGS) github.com/mdm-code/tqweb/...

cover:
	@$(GO) test -coverprofile=$(COV_PROFILE) -covermode=atomic ./...
	@$(GO) tool cover -html=$(COV_PROFILE)

clean:
	@$(GO) clean github.com/mdm-code/tqweb/...
	@$(GO) mod tidy
	@$(GO) clean -testcache
	@rm -f $(COV_PROFILE)

update-js:
	@./tools/scripts/dl-htmx

gen-templ:
	@templ generate

serve:
	@$(GO) build -C ./cmd/tqweb -trimpath -o tqweb
	@./cmd/tqweb/tqweb
