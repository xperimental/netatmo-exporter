SHELL := /bin/bash
GO ?= go
GO_CMD := CGO_ENABLED=0 $(GO)
GIT_VERSION := $(shell git describe --tags --dirty)
VERSION := $(GIT_VERSION:v%=%)
GIT_COMMIT := $(shell git rev-parse HEAD)
GITHUB_REF ?= refs/heads/master
DOCKER_TAG != if [[ "$(GITHUB_REF)" == "refs/heads/master" ]]; then \
		echo "latest"; \
	else \
		echo "$(VERSION)"; \
	fi

.PHONY: all
all: test build-binary

.PHONY: test
test:
	$(GO_CMD) test -cover ./...

.PHONY: build-binary
build-binary:
	$(GO_CMD) build -tags netgo -ldflags "-w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o netatmo-exporter .

.PHONY: image
image:
	docker build -t "xperimental/netatmo-exporter:$(DOCKER_TAG)" .

.PHONY: all-images
all-images:
	docker buildx build -t "ghcr.io/xperimental/netatmo-exporter:$(DOCKER_TAG)" -t "xperimental/netatmo-exporter:$(DOCKER_TAG)" --platform linux/amd64,linux/arm64 --push .

.PHONY: clean
clean:
	rm -f netatmo-exporter
