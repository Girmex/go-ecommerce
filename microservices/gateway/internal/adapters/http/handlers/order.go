package handlers

import (
	"encoding/json"
	"net/http"

	ordergrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/order"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
	orderproto "github.com/Girmex/go-ecommerce/microservices/order/proto"

	"google.golang.org/grpc/metadata"
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
