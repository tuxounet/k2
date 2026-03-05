# ─────────────────────────────────────────────────────────────
#  k2 — A declarative YAML-driven template engine
#
#  k2 lets you define reusable templates, wire them into an
#  inventory, and generate entire project trees with a single
#  command.  Plan → Apply → Destroy — simple as that.
#
#  Run `make` or `make help` to see available targets.
# ─────────────────────────────────────────────────────────────

APP_NAME     := k2
GIT_BRANCH   := $(shell git rev-parse --abbrev-ref HEAD)
GO_PATH      := $(shell go env GOPATH)
VERSION      := $(if $(CI_COMMIT_TAG),$(CI_COMMIT_TAG),v$(GIT_BRANCH))
VERSION_FILE := ./version.txt

# ── Default target ───────────────────────────────────────────
.DEFAULT_GOAL := help

help: ## Show this help
	@echo ""
	@echo "  k2 — declarative YAML-driven template engine"
	@echo ""
	@echo "  Usage:  make <target>"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ── Build & Run ──────────────────────────────────────────────
write-version: ## Write current version to version.txt
	@echo $(VERSION) > $(VERSION_FILE)

run: write-version ## Run k2 locally (go run)
	go run ./main.go

build: write-version ## Build the k2 binary into .out/
	go build -o ./.out/k2 ./main.go

bump-patch: build ## Bump the patch version number
	@echo "Bumping version patch"
	@bash ./tools/bump-patch.sh

# ── Quality ──────────────────────────────────────────────────
test: ## Run all tests
	go test -v ./...

coverage: ## Run tests with coverage report
	go test -v -coverprofile=coverage.out ./...

# ── Samples (quick try with the bundled inventory) ───────────
plan: ## Plan changes using the sample inventory
	go run ./main.go plan --inventory ./samples/k2.inventory.yaml

apply: ## Apply templates using the sample inventory
	go run ./main.go apply --inventory ./samples/k2.inventory.yaml

destroy: ## Destroy generated files using the sample inventory
	go run ./main.go destroy --inventory ./samples/k2.inventory.yaml

.PHONY: help write-version run build bump-patch test coverage plan apply destroy