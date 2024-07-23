ROOT_PATH := .
BINARY_NAME := probe-mux

.PHONY: deps
deps:
	@go fmt .
	@go mod tidy -v

.PHONY: build
build: deps
	@go build -o=${ROOT_PATH}/${BINARY_NAME} ${ROOT_PATH}

.PHONY: run
run: build
	@-./${BINARY_NAME} $(ARGS)
	@$(MAKE) clean

.PHONY: test
test: build
	@go test -v ./...
	@$(MAKE) clean

.PHONY: clean
clean:
	@rm -rf ./${BINARY_NAME}