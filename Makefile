include common.mk
include tools.mk

LDFLAGS += -X "$(MODULE)/version.Version=$(VERSION)" -X "$(MODULE)/version.CommitSHA=$(VERSION_HASH)"

## Build:

.PHONY: build
build: | build-frontend build-backend ## Build binary

.PHONY: build-frontend
build-frontend: ## Build frontend
	$Q cd frontend && npm ci && npm run build

.PHONY: build-backend
build-backend: ## Build backend
	$Q $(go) build -ldflags '$(LDFLAGS)' -o .

.PHONY: test
test: | test-frontend test-backend ## Run all tests

.PHONY: test-frontend
test-frontend: ## Run frontend tests

.PHONY: test-backend
test-backend: ## Run backend tests
	$Q $(go) test -v ./...

.PHONY: lint
lint: lint-frontend lint-backend ## Run all linters

.PHONY: lint-frontend
lint-frontend: ## Run frontend linters
	$Q cd frontend && npm ci && npm run lint

.PHONY: lint-backend
lint-backend: | $(golangci-lint) ## Run backend linters
    ## https://github.com/golangci/golangci-lint/issues/825
	$Q go mod tidy && go mod vendor && $(golangci-lint) run -v

fmt: $(goimports) ## Format source files
	$Q $(goimports) -local $(MODULE) -w $$(find . -type f -name '*.go' -not -path "./vendor/*")

clean: clean-tools ## Clean

## Help:
help: ## Show this help
	@echo ''
	@echo 'Usage:'
	@echo '  ${YELLOW}make${RESET} ${GREEN}<target> [options]${RESET}'
	@echo ''
	@echo 'Options:'
	@$(call global_option, "V [0|1]", "enable verbose mode (default:0)")
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} { \
		if (/^[a-zA-Z_-]+:.*?##.*$$/) {printf "    ${YELLOW}%-20s${GREEN}%s${RESET}\n", $$1, $$2} \
		else if (/^## .*$$/) {printf "  ${CYAN}%s${RESET}\n", substr($$1,4)} \
		}' $(MAKEFILE_LIST)
