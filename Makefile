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
.PHONY: all test build exe run clean install
