GIT_VERSION ?= $(shell git describe --abbrev=4 --dirty --always --tags)

GOPATH ?= $(HOME)/go
BIN_DIR = $(GOPATH)/bin
TMPDIR ?= $(shell dirname $$(mktemp -u))


# Project specific variables

PACKAGE = doc_bot
APP_NAME ?= $(PACKAGE)
NAMESPACE = github.com/AleksandrAkhapkin/$(PACKAGE)
COVER_FILE ?= $(TMPDIR)/$(PACKAGE)-coverage.out

CONFIG_LOGGER_LEVEL ?= debug

# PostgreSQL
DB_USER ?= docbot
DB_PASSWORD ?= secret
DB_NAME ?= docbot
DB_HOST ?= localhost
DB_PORT ?= 5432
DB_SSL_MODE ?= disable
DB_SCHEMA ?= doc_bot
DATABASE_URL_WO_SCHEMA ?= postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=$(DB_SSL_MODE)
DATABASE_URL ?= $(DATABASE_URL_WO_SCHEMA)&search_path=$(DB_SCHEMA)

TEST_DB_NAME ?= docbot-test-db
TEST_DATABASE_URL ?= postgresql://$(DB_USER):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(TEST_DB_NAME)?sslmode=$(DB_SSL_MODE)&search_path=$(DB_SCHEMA)

.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@echo Installing tools from tools.go
	@grep '_ "' tools.go | grep -o '"[^"]*"' | xargs -tI % go install %

# Main targets

all: test build
.DEFAULT_GOAL := all

.PHONY: build
build: ## Build the project binary
	go build -ldflags "-X main.version=$(GIT_VERSION)" ./cmd/$(PACKAGE)/

.PHONY: test
test: ## Run unit (short) tests
	go test -short ./... -coverprofile=$(COVER_FILE)
	go tool cover -func=$(COVER_FILE) | grep ^total

$(COVER_FILE):
	$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	go tool cover -html=$(COVER_FILE)
	rm -f $(COVER_FILE)

.PHONY: test_integration
test_integration: ## Run integration tests, creates new test db and can break your data
	echo "CREATE DATABASE \"$(TEST_DB_NAME)\";" | psql "$(DATABASE_URL_WO_SCHEMA)"
	DB_NAME=$(TEST_DB_NAME) $(MAKE) migrate
	TEST_DATABASE_URL="$(TEST_DATABASE_URL)" go test ./... -tags=integration

.PHONY: bench
bench: ## Run benchmarks
	go test ./... -short -bench=. -run="Benchmark*"

.PHONY: lint
lint: tools ## Check the project with lint
	golint -set_exit_status ./...

.PHONY: vet
vet: ## Check the project with vet
	go vet ./...

.PHONY: fmt
fmt: ## Run go fmt for the whole project
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do gofmt -e -l -w $$d/*.go; done)

.PHONY: imports
imports: tools ## Check and fix import section by import rules
	test -z $$(for d in $$(go list -f {{.Dir}} ./...); do goimports -e -l -local $(NAMESPACE) -w $$d/*.go; done)

.PHONY: code_style
code_style: ## Check code style issues in the project, only line length at the moment
	find . -name '*.go' ! -name '*_gen.go' -not -path "./.*" | grep -v _test.go | grep -v packed-packr.go | grep -v docs.go | xargs -i sh -c "expand -t 4 {} | awk 'length>120' | grep -v '//' && expand -t 4 {} | awk 'length>120' | grep -v '@Param' | grep -v '//' |  grep -v 'json:' | grep -v 'gorm:' | wc -l | grep '^0$$' > /dev/null 2>&1 || echo {}"
	find . -name '*.go' ! -name '*_gen.go' -not -path "./.*" | grep -v _test.go | grep -v packed-packr.go | grep -v docs.go | xargs expand -t 4 | awk 'length>120' | grep -v '//' | grep -v '@Param' | grep -v 'json:' | grep -v 'gorm:' | wc -l | grep '^0$$'

.PHONY: static_check
static_check: fmt imports vet lint code_style ## Run static checks (fmt, lint, imports, vet, ...) all over the project

.PHONY: check
check: static_check test ## Check project with static checks and unit tests

.PHONY: create_schema
create_schema: ## Create schema in database
	echo "CREATE SCHEMA IF NOT EXISTS \"$(DB_SCHEMA)\";" | psql "$(DATABASE_URL_WO_SCHEMA)"

.PHONY: migrate
migrate: build ## Run migrations
	APP_NAME="$(APP_NAME)" \
	CONFIG_DATABASE_ARGS="$(DATABASE_URL)" \
	./$(APP_NAME) -migrate

.PHONY: migrate_dev
migrate_dev: create_schema migrate ## Run migrations with create schema, dev only

.PHONY: dependencies
dependencies: ## Manage go mod dependencies, beautify go.mod and go.sum files
	go mod tidy

.PHONY: run
run: build ## Start the project
	APP_NAME="$(APP_NAME)" \
	./$(PACKAGE)

.PHONY: clean
clean: ## Clean the project from built files
	rm -f ./$(PACKAGE) $(COVER_FILE)

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: docs
docs: tools
	swag init -g cmd/${PACKAGE}/main.go --output "./docs" --parseDependency --parseInternal

.PHONY: dependencies-download
dependencies-download:
	go mod download
