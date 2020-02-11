PROJECT_NAME := "terraform-provider-hiera5"
PKG := "gitlab.com/sbitio/$(PROJECT_NAME)"
PKG_LIST := $(shell go list ${PKG}/... | grep -v /vendor/)
GO_FILES := $(shell find . -name '*.go' | grep -v /vendor/ | grep -v _test.go)

GOFMT_FILES?=$$(find . -name '*.go' |grep -v vendor)
TEST?=$$(go list ./... |grep -v 'vendor')

.PHONY: help all build clean coverage coverhtml dep errcheck fmt fmtcheck lint test testtf testacc vet test-compile

default: build

all: build

build: dep ## Build the binary file
	@go build -v -o $(PROJECT_NAME)

clean: ## Remove previous build
	@rm -f $(PROJECT_NAME)

coverage: ## Generate global code coverage report
	TF_ACC=1 ./scripts/coverage.sh;

coverhtml: ## Generate global code coverage report in HTML
	TF_ACC=1 ./scripts/coverage.sh html;

dep: ## Get the dependencies
        ifneq (,$(findstring "-mod=vendor",$(GOFLAGS)))
	  go get -v -d ./...
        endif

errcheck:
	@sh -c "'$(CURDIR)/scripts/errcheck.sh'"

fmt:
	gofmt -w $(GOFMT_FILES)

fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint: ## Lint the files
	@golint -set_exit_status ${PKG_LIST}

msan: dep ## Run memory sanitizer
	CC=clang go test -msan -short ${PKG_LIST}

race: dep ## Run data race detector
	@go test -race -short ${PKG_LIST}

test: ## Run unittests
	@go test -short ${PKG_LIST}

testacc: fmtcheck
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

testtf: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

test-compile:
	@if [ "$(TEST)" = "./..." ]; then \
		echo "ERROR: Set TEST to a specific package. For example,"; \
		echo "  make test-compile TEST=./$(PKG_NAME)"; \
		exit 1; \
	fi
	go test -c $(TEST) $(TESTARGS)

vet:
	@echo "go vet ."
	@go vet $$(go list ./... | grep -v vendor/) ; if [ $$? -eq 1 ]; then \
	  	echo ""; \
		echo "Vet found suspicious constructs. Please check the reported constructs"; \
		echo "and fix them if necessary before submitting the code for review."; \
		exit 1; \
	fi

help: ## Display this help screen
	@grep -h -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
