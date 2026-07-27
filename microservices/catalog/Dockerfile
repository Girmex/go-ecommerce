# ---------- Builder ----------
FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux \
    go build -o catalog-service ./microservices/catalog/cmd

# ---------- Runtime ----------
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/catalog-service .

EXPOSE 50052

CMD ["./catalog-service"]