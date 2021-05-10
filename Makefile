TARGET = model
COVERAGE_TARGET = count.out

BIN_DIR = bin
CMD_DIR = cmd
INT_DIR = internal
PKG_DIR = pkg
COVERAGE_DIR = coverage

INT_PACKAGES = $(shell find $(INT_DIR) -depth 1)
PUB_PACKAGES = $(shell find $(PKG_DIR) -depth 1)

.PHONY: all
all: build

.PHONY: build
build: env dep lint test $(BIN_DIR)
	go build -o $(BIN_DIR)/$(TARGET) $(CMD_DIR)/*

.PHONY: lint
lint:
	golangci-lint run

.PHONY: dep
dep:
	dep ensure

.PHONY: test
test: $(COVERAGE_DIR)
	go test -v -covermode=count -coverprofile=$(COVERAGE_DIR)/$(COVERAGE_TARGET) ./...

.PHONY: env
env:
	go env -w GO111MODULE=auto

.PHONY: fmt
fmt:
	gofmt -s -w .
	goimports -w .

$(COVERAGE_DIR):
	mkdir -p $(COVERAGE_DIR)
	touch $(COVERAGE_DIR)/$(COVERAGE_TARGET)

$(BIN_DIR):
	mkdir -p $(BIN_DIR)

.PHONY: clean
clean:
	rm -rf $(COVERAGE_DIR) $(BIN_DIR)
