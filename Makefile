export VERSION ?= $(shell cat ./VERSION)
export CLI_BIN ?= dist/difff_linux_amd64_v1/difff


.PHONY: build
build:
	goreleaser release --snapshot --clean

.PHONY: release
release:
	goreleaser release --clean

.PHONY: run
run: build
	$(CLI_BIN) \
		e2e/data/source \
		e2e/data/target

.PHONY: help
help: build
	$(CLI_BIN) --help

.PHONY: test
test:
	go test -v -cover -coverprofile=index.out ./...

.PHONY: cover
cover: test
	go tool cover -html=index.out -o index.html
	python3 -m http.server 8765

.PHONY: e2e
e2e:
	./e2e/test.sh

.PHONY: golangci
golangci:
	golangci-lint run -v ./...

.PHONY: govulncheck
govulncheck:
	govulncheck ./...

.PHONY: ci
ci: build test e2e golangci govulncheck

.PHONY: tag
tag:
	git tag -f $(VERSION)
	git push origin $(VERSION)
