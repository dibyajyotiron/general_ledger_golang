.PHONY: build clean tool lint help

all: build

build:
	@go build -v .

tool:
	go vet general_ledger_golang/...; true
	gofmt -w .

lint:
	golint general_ledger_golang/...

clean:
	rm -rf go-gin-example
	go clean -i .

help:
	@echo "make: compile packages and dependencies"
	@echo "make tool: run specified go tool"
	@echo "make lint: golint general_ledger_golang/..."
	@echo "make clean: remove object files and cached files"
