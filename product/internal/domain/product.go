package domain

import (
	"errors"
	"time"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrInsufficientStock = errors.New("insufficient stock")
)

type Product struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ReserveStock decrements stock for an order, guarding against overselling.
func (p *Product) ReserveStock(qty int) error {
	if p.Stock < qty {
		return ErrInsufficientStock
	}
	p.Stock -= qty
	return nil
}
