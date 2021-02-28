REPO_PATH = ~/Github/distr-model

.PHONY: build
	build: bin/distr-model

bin/search-ctl:
	@echo "===> $@"
	@go build -i -o bin/distr-model -v $(REPO_PATH)/priorityq

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@echo "===> $@"
	@go test ${GOTESTFLAGS} `go list ./...`

.PHONY: clean
clean:
	@echo "===> $@"
	@rm -rf bin