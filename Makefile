.PHONY: \
	proto-catalog \
	run-catalog \
	build \
	test

proto-catalog:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/catalog/api/proto/catalog.proto

run-catalog:
	APP_ENV=dev go run ./microservices/catalog/cmd

build:
	go build ./microservices/catalog/...

test:
	go test ./microservices/catalog/...