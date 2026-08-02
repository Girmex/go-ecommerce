package domain

import "time"

type Payment struct {
	ID        uint
	OrderID   uint
	UserID    uint
	Amount    float64
	Status    PaymentStatus
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (p *Payment) Complete() error {

	if p.Status == PaymentSuccess {
		return ErrPaymentAlreadyCompleted
	}

	if p.Status == PaymentFailed {
		return ErrPaymentAlreadyFailed
	}

	p.Status = PaymentSuccess

	return nil
}

func (p *Payment) Fail() error {

	if p.Status == PaymentSuccess {
		return ErrPaymentAlreadyCompleted
	}

	if p.Status == PaymentFailed {
		return ErrPaymentAlreadyFailed
	}

	p.Status = PaymentFailed

	return nil
}
