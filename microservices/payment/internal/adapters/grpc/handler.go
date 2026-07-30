package grpc

import (
	"context"

	"github.com/Girmex/go-ecommerce/microservices/payment/internal/application"
	proto "github.com/Girmex/go-ecommerce/microservices/payment/proto"
)

type Handler struct {
	proto.UnimplementedPaymentServiceServer

	service *application.PaymentService
}

func NewHandler(
	service *application.PaymentService,
) *Handler {

	return &Handler{
		service: service,
	}
}

func (h *Handler) CreatePayment(
	ctx context.Context,
	req *proto.CreatePaymentRequest,
) (*proto.Payment, error) {

	payment, err := h.service.CreatePayment(
		ctx,
		uint(req.OrderId),
		uint(req.UserId),
		req.Amount,
	)
	if err != nil {
		return nil, toStatusError(err)
	}

	return toProtoPayment(payment), nil
}