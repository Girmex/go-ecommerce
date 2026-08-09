.PHONY: \
	proto-auth \
	proto-catalog \
	proto-order \
	proto-payment \
	run-auth \
	run-catalog \
	run-order \
	run-payment \
	run-gateway \
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
	go run ./microservices/auth/cmd


run-catalog:
	go run ./microservices/catalog/cmd

run-order:
	go run ./microservices/order/cmd

run-payment:
	go run ./microservices/payment/cmd

run-notification:
	go run ./microservices/notification/cmd

run-gateway:
	go run ./microservices/gateway/cmd
test:
	go test ./...