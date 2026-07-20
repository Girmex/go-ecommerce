# ---------- Builder ----------
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/microservices/catalog

RUN CGO_ENABLED=0 GOOS=linux go build -o catalog-service ./cmd

# ---------- Runtime ----------
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/microservices/catalog/catalog-service .

EXPOSE 50051

CMD ["./catalog-service"]