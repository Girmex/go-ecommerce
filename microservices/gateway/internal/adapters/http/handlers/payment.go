package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	paymentgrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/payment"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
	paymentproto "github.com/Girmex/go-ecommerce/microservices/payment/proto"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/metadata"
)

type PaymentHandler struct {
	client *paymentgrpc.Client
}

func NewPaymentHandler(client *paymentgrpc.Client) *PaymentHandler {
	return &PaymentHandler{
		client: client,
	}
}
func (h *PaymentHandler) CreatePayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	grpcReq := &paymentproto.CreatePaymentRequest{
		OrderId: req.OrderID,
		UserId:  uint32(claims.UserID),
		Amount:  req.Amount,
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.CreatePayment(
		ctx,
		grpcReq,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) GetPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid payment id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.GetPayment(
		ctx,
		&paymentproto.GetPaymentRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *PaymentHandler) CompletePayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid payment id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.CompletePayment(
		ctx,
		&paymentproto.CompletePaymentRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (h *PaymentHandler) FailPayment(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid payment id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.FailPayment(
		ctx,
		&paymentproto.FailPaymentRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
