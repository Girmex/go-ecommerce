# Go E-Commerce Microservices

A simple e-commerce backend built with **Go** using a microservices architecture.

## Services

- **User Service** — authentication and user management
- **Product Service** — product and stock management
- **Order Service** — order management
- **Payment Service** — payment processing

## Tech Stack

- Go
- Chi
- PostgreSQL
- JWT
- Swagger
- Docker
- Hexagonal Architecture

## Structure

```text id="m7z0nc"
chi-ecommerce-ms/
├── user/
├── product/
├── order/
├── payment/
└── README.md
```

Each service has its own:

- `go.mod`
- Docker Compose configuration
- Database migrations
- Swagger documentation
- Configuration

Each service can be run independently from its own directory.]633;E;echo "";2cb93968-fbf9-4130-814b-34c329917fe8]633;C
## Swagger Testing

### password-validation


### login-response


### list-products


### get-product


### post-product


### place-order


### get-order


### cancel-order


### charge-payment


### get-payment


### list-payment


### refund-payment


### protected-endpoint


]633;E;echo "";2cb93968-fbf9-4130-814b-34c329917fe8]633;C
## Swagger Testing

### password-validation

![1.password-validation](screenshots/1.password-validation.png)

### login-response

![2.login-response](screenshots/2.login-response.png)

### list-products

![3.list-products](screenshots/3.list-products.png)

### get-product

![4.get-product](screenshots/4.get-product.png)

### post-product

![5.post-product](screenshots/5.post-product.png)

### place-order

![6.place-order](screenshots/6.place-order.png)

### get-order

![7.get-order](screenshots/7.get-order.png)

### cancel-order

![8.cancel-order](screenshots/8.cancel-order.png)

### charge-payment

![9.charge-payment](screenshots/9.charge-payment.png)

### get-payment

![10.get-payment](screenshots/10.get-payment.png)

### list-payment

![11.list-payment](screenshots/11.list-payment.png)

### refund-payment

![12.refund-payment](screenshots/12.refund-payment.png)

### protected-endpoint

![protected-endpoint](screenshots/protected-endpoint.png)

