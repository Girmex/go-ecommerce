.PHONY: \
	proto-auth \
	proto-catalog \
	proto-order \
	proto-payment \
	run-auth \
	run-catalog \
	run-order \
	run-payment \
	test

proto-auth:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/auth/proto/auth.proto
proto-catalog:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/catalog/proto/catalog.proto

proto-order:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/order/proto/order.proto

proto-payment:
	protoc \
		--proto_path=. \
		--go_out=. \
		--go_opt=paths=source_relative \
		--go-grpc_out=. \
		--go-grpc_opt=paths=source_relative \
		microservices/payment/proto/payment.proto	

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

run-order:
	APP_NAME="Order Service" \
	APP_ENV=development \
	GRPC_PORT=50053 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=root \
	POSTGRES_PASSWORD=root \
	POSTGRES_DB=order_db \
	JWT_SECRET="super-secret-key" \
	go run ./microservices/order/cmd

run-payment:
	APP_NAME="Payment Service" \
	APP_ENV=development \
	GRPC_PORT=50054 \
	POSTGRES_HOST=localhost \
	POSTGRES_PORT=5432 \
	POSTGRES_USER=root \
	POSTGRES_PASSWORD=root \
	POSTGRES_DB=payment_db \
	JWT_SECRET="super-secret-key" \
	go run ./microservices/payment/cmd
test:
	go test ./...