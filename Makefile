VERSION = 0.1.0
TAG = $(VERSION)
PREFIX = nginx-exporter


GIT_COMMIT = $(shell git rev-parse HEAD)
DATE = $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")


.PHONY: nginx-exporter
nginx-exporter:
	CGO_ENABLED=0 go build -trimpath -ldflags "-s -w -X main.version=$(VERSION) -X main.commit=$(GIT_COMMIT) -X main.date=$(DATE)" -o nginx-exporter

.PHONY: build-goreleaser
build-goreleaser: ## Build all binaries using GoReleaser
	GOPATH=$(shell go env GOPATH) goreleaser build --rm-dist --snapshot

.PHONY: deps
deps:
	@go mod tidy && go mod verify && go mod download

.PHONY: clean
clean:
	-rm -r dist
	-rm nginx-exporter
