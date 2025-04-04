ENV ?= stg
GO_CMD ?= go

run:
	$(GO_CMD) run cmd/app/main.go --env=$(ENV)

test:
	$(GO_CMD) test ./... -v

fmt:
	$(GO_CMD) fmt ./...

clean:
	$(GO_CMD) clean -testcache

lint:
	golangci-lint run ./...

ci: fmt lint test
