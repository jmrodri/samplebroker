REGISTRY         ?= docker.io
ORG              ?= jmrodri
TAG              ?= latest
BROKER_IMAGE     ?= $(REGISTRY)/$(ORG)/samplebroker
BUILD_DIR        = "${GOPATH}/src/github.com/jmrodri/samplebroker/build"
SOURCE_DIRS      = cmd pkg
SOURCES          := $(shell find . -name '*.go' -not -path "*/vendor/*")
PACKAGES         := $(shell go list ./pkg/...)
COVERAGE_SVC     := travis-ci
.DEFAULT_GOAL    := build

vendor: ## Install or update project dependencies
	@dep ensure

samplebroker: $(SOURCES) ## Build the samplebroker
	go build -i -ldflags="-s -w" ./cmd/samplebroker

build: samplebroker ## Build binary from source
	@echo > /dev/null

lint: ## Run golint
	@golint -set_exit_status $(addsuffix /... , $(SOURCE_DIRS))

fmtcheck: ## Check go formatting
	@gofmt -l $(SOURCES) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

test: ## Run unit tests
	@go test -cover ./pkg/...

coverage-all.out: $(SOURCES)
	@echo "mode: count" > coverage-all.out
	@$(foreach pkg,$(PACKAGES),\
		go test -coverprofile=coverage.out -covermode=count $(pkg);\
		tail -n +2 coverage.out >> coverage-all.out;)

test-coverage-html: coverage-all.out ## checkout the coverage locally of your tests
	@go tool cover -html=coverage-all.out

vet: ## Run go vet
	@go tool vet ./cmd ./pkg

check: fmtcheck vet lint build test ## Pre-flight checks before creating PR

build-image: ## Build a docker image with the broker binary
	env GOOS=linux go build -i -ldflags="-s -s" -o ${BUILD_DIR}/samplebroker ./cmd/samplebroker
	docker build -f ${BUILD_DIR}/Dockerfile -t ${BROKER_IMAGE}:${TAG} ${BUILD_DIR}
	@echo
	@echo "Remember you need to push your image before calling make deploy"
	@echo "    docker push ${BROKER_IMAGE}:${TAG}"

clean: ## Clean up your working environment
	@rm -f samplebroker
	@rm -f build/samplebroker
	@rm -f coverage-all.out coverage.out handler.out registries.out validation.out

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

.PHONY: build-image clean lint build fmtcheck test vet help test-coverage-html
