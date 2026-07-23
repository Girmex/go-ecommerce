.PHONY: \
	proto-auth \
	run-auth \
	build \
	test


proto-auth:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/auth/api/proto/auth.proto

run-auth:
	APP_NAME="Auth Service" \
	APP_ENV=development \
	GRPC_PORT=50052 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5434 \
	POSTGRES_USER=root \
	POSTGRES_PASSWORD=root \
	POSTGRES_DB=auth_db \
	go run ./microservices/auth/cmd
build:
	go build ./microservices/auth/...

test:
	go test ./microservices/auth/...