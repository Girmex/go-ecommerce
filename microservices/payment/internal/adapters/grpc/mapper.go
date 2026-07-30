package grpc

import (
	"github.com/Girmex/go-ecommerce/microservices/payment/internal/domain"
	proto "github.com/Girmex/go-ecommerce/microservices/payment/proto"
)

func toProtoStatus(status domain.PaymentStatus) proto.PaymentStatus {

	switch status {

	case domain.PaymentPending:
		return proto.PaymentStatus_PAYMENT_STATUS_PENDING

	case domain.PaymentSuccess:
		return proto.PaymentStatus_PAYMENT_STATUS_SUCCESS

	case domain.PaymentFailed:
		return proto.PaymentStatus_PAYMENT_STATUS_FAILED

	default:
		return proto.PaymentStatus_PAYMENT_STATUS_UNSPECIFIED
	}
}

func toProtoPayment(payment *domain.Payment) *proto.Payment {

	return &proto.Payment{
		Id:      uint32(payment.ID),
		OrderId: uint32(payment.OrderID),
		UserId:  uint32(payment.UserID),
		Amount:  payment.Amount,
		Status:  toProtoStatus(payment.Status),
	}
}
