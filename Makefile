.PHONY: run test fmt

run:
	cd user && go run ./cmd

test:
	cd user && go test ./...

fmt:
	cd user && gofmt -w .
