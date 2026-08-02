package ports

import "context"

type OrderClient interface {
	MarkOrderAsPaid(
		ctx context.Context,
		orderID uint,
	) error

	CancelOrder(
		ctx context.Context,
		orderID uint,
	) error
}
