.PHONY: all test build-binary image clean

GO ?= go
GO_OS ?= linux
GO_ARCH ?= amd64
GO_CMD := CGO_ENABLED=0 $(GO)
GIT_VERSION := $(shell git describe --tags --dirty)
VERSION := $(GIT_VERSION:v%=%)
GIT_COMMIT := $(shell git rev-parse HEAD)

all: test build-binary

test:
	$(GO_CMD) test -cover ./...

build-binary:
	GOOS=$(GO_OS) GOARCH=$(GO_ARCH) $(GO_CMD) build -tags netgo -ldflags "-w -X main.Version=$(VERSION) -X main.GitCommit=$(GIT_COMMIT)" -o netatmo-exporter .

image:
	docker build -t "xperimental/netatmo-exporter:$(VERSION)" .

clean:
	rm -f netatmo-exporter
