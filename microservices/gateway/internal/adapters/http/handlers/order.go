package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	ordergrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/order"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
	orderproto "github.com/Girmex/go-ecommerce/microservices/order/proto"
	"github.com/go-chi/chi/v5"

	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type OrderHandler struct {
	client *ordergrpc.Client
}

func NewOrderHandler(
	client *ordergrpc.Client,
) *OrderHandler {

	return &OrderHandler{
		client: client,
	}
}

func (h *OrderHandler) CreateOrder(
	w http.ResponseWriter,
	r *http.Request,
) {

	var req dto.CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusBadRequest,
		)
		return
	}

	claims, ok := middleware.Claims(r.Context())

	if !ok {
		http.Error(
			w,
			"unauthorized",
			http.StatusUnauthorized,
		)
		return
	}

	items := make(
		[]*orderproto.OrderItemRequest,
		0,
		len(req.Items),
	)

	for _, item := range req.Items {

		items = append(
			items,
			&orderproto.OrderItemRequest{
				ProductId: item.ProductID,
				Quantity:  item.Quantity,
			},
		)
	}

	_ = claims

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.CreateOrder(
		ctx,
		&orderproto.CreateOrderRequest{
			Items: items,
		},
	)

	if err != nil {
		http.Error(
			w,
			err.Error(),
			http.StatusInternalServerError,
		)
		return
	}

	w.Header().Set(
		"Content-Type",
		"application/json",
	)

	json.NewEncoder(w).Encode(resp)
}

func (h *OrderHandler) ListOrders(
	w http.ResponseWriter,
	r *http.Request,
) {
	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.ListOrders(
		ctx,
		&emptypb.Empty{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
func (h *OrderHandler) GetOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.GetOrder(
		ctx,
		&orderproto.GetOrderRequest{
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
func (h *OrderHandler) CancelOrder(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	_, err = h.client.CancelOrder(
		ctx,
		&orderproto.GetOrderRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}


func (h *OrderHandler) UpdateOrderStatus(
	w http.ResponseWriter,
	r *http.Request,
) {
	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid order id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateOrderStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.UpdateOrderStatus(
		ctx,
		&orderproto.UpdateOrderStatusRequest{
			Id:     uint32(id),
			Status: orderproto.OrderStatus(req.Status),
		},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
