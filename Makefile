help:
	@echo "make [all|test|build|exe|run|clean|install]"
	@echo "all: test build exe run"
	@echo "test: run unit test"
	@echo "build: build the module"
	@echo "exe: make executable for the module"
	@echo "clean: clean module C objects"
	@echo "run: exec the module code"
	@echo "install: install the module in go libs"
all: test build exe run
	@echo "Make all scopes"
build:
	@go build . > /dev/null
exe:
	@go build --buildmode exe .
run:
	@go run main.go
install:
	@go install
clean:
	@go clean
test:
	@go test
.PHONY: help all test build exe run clean install
