INT_DIR = internal
PKG_DIR = pkg

INT_PACKAGES = $(shell find $(INT_DIR) -depth 1)
PUB_PACKAGES = $(shell find $(PKG_DIR) -depth 1)

.PHONY: all
all: env test

.PHONY: test
test:
	go test -v -cover ./...

.PHONY: env
env:
	go env -w GO111MODULE=auto

.PHONY: fmt
fmt:
	gofmt -s -w .
