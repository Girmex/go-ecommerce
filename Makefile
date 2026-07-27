.PHONY: \
	proto-auth \
	run-auth \
	run-catalog \
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
	GRPC_PORT=50051 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=root \
	POSTGRES_PASSWORD=root \
	POSTGRES_DB=auth_db \
	JWT_SECRET="super-secret-key" \
	go run ./microservices/auth/cmd

run-catalog:
	APP_NAME="Catalog Service" \
	APP_ENV=development \
	GRPC_PORT=50052 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=root \
	POSTGRES_PASSWORD=root \
	POSTGRES_DB=catalog_db \
	JWT_SECRET="super-secret-key" \
	go run ./microservices/catalog/cmd

test:
	go test ./microservices/auth/...