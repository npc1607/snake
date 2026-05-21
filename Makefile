APP := snake
CMD := ./cmd/$(APP)
BIN_DIR := bin
BIN := $(BIN_DIR)/$(APP)

.DEFAULT_GOAL := help

.PHONY: help
help: ## Show available targets.
	@awk 'BEGIN {FS = ":.*##"; printf "Usage:\n  make <target>\n\nTargets:\n"} /^[a-zA-Z_-]+:.*##/ {printf "  %-12s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.PHONY: run
run: ## Run the game.
	go run $(CMD)

.PHONY: build
build: ## Build the game binary into bin/.
	@mkdir -p $(BIN_DIR)
	go build -o $(BIN) $(CMD)

.PHONY: test
test: ## Run all tests.
	go test ./...

.PHONY: fmt
fmt: ## Format Go source files.
	gofmt -w $$(find . -name '*.go' -not -path './vendor/*')

.PHONY: tidy
tidy: ## Tidy Go module dependencies.
	go mod tidy

.PHONY: clean
clean: ## Remove build artifacts.
	rm -rf $(BIN_DIR)
