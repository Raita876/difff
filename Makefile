VERSION := $(shell cat ./VERSION)
CLI_NAME := difff

export GOOS ?= linux
export GOARCH ?= amd64

.PHONY: build
build:
	go build \
		-o ./bin/$(GOOS)/$(GOARCH)/$(CLI_NAME) \
		-ldflags "-X main.version=$(VERSION) -X main.name=$(CLI_NAME)" \
		./cmd/$(CLI_NAME)

.PHONY: run
run: build
	./bin/$(GOOS)/$(GOARCH)/$(CLI_NAME) \
		e2e/data/source \
		e2e/data/target

.PHONY: help
help: build
	./bin/$(GOOS)/$(GOARCH)/$(CLI_NAME) --help

.PHONY: test
test:
	go test -v -cover -coverprofile=index.out ./...

.PHONY: cover
cover: test
	go tool cover -html=index.out -o index.html
	python3 -m http.server 8765

.PHONY: golangci
golangci:
	golangci-lint run -v ./...

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: ci
ci: test golangci govulncheck
