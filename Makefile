# Set default shell to bash
SHELL := /bin/bash -o pipefail -o errexit -o nounset

COVERAGE_REPORT ?= coverage.txt

## Format go code
.PHONY: fmt
fmt:
	go run golang.org/x/tools/cmd/goimports@v0.1.7 -w .

## lint code
.PHONY: lint
lint:
	go run github.com/golangci/golangci-lint/cmd/golangci-lint@v1.43.0 run ./...

## generates coverage report
.PHONY: coverage
coverage: 
	go test -coverprofile=$(COVERAGE_REPORT) ./...

## generates coverage report and shows it on the browser locally
.PHONY: coverage/show
coverage/show: coverage
	go tool cover -html=$(COVERAGE_REPORT)

## test code
.PHONY: test
test: 
	go test -race ./...

install:
	go install ./cmd/mdtoc
