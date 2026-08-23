package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/adapter/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/order/internal/port"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type OrderHandler struct {
	orderService port.OrderService
	validator    *validator.Validate
}

func NewOrderHandler(orderService port.OrderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
		validator:    validator.New(),
	}
}

// @Summary      Place an order
// @Description  Creates an order and reserves stock for each product.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        order  body      CreateOrderRequest  true  "Order payload"
// @Success      201    {object}  OrderResponse
// @Failure      400    {object}  ErrorResponse
// @Failure      401    {object}  ErrorResponse
// @Failure      404    {object}  ErrorResponse
// @Failure      500    {object}  ErrorResponse
// @Router       /orders [post]
func (h *OrderHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateOrderRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_BODY",
			"request body is not valid JSON",
			nil,
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"request failed validation",
			err.Error(),
		)
		return
	}

	items := make([]port.CreateOrderItemInput, 0, len(req.Items))

	for _, item := range req.Items {
		items = append(items, port.CreateOrderItemInput{
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
		})
	}
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
			nil,
		)
		return
	}

	order, err := h.orderService.PlaceOrder(
		r.Context(),
		port.CreateOrderInput{
			UserID:        userID,
			Items:         items,
			PaymentMethod: req.PaymentMethod,
		},
	)

	if err != nil {
		switch {
		case errors.Is(err, domain.ErrEmptyOrder):
			writeError(
				w,
				http.StatusBadRequest,
				"EMPTY_ORDER",
				err.Error(),
				nil,
			)

		case errors.Is(err, domain.ErrUserNotFound):
			writeError(
				w,
				http.StatusNotFound,
				"USER_NOT_FOUND",
				err.Error(),
				nil,
			)

		case errors.Is(err, domain.ErrProductNotFound):
			writeError(
				w,
				http.StatusNotFound,
				"PRODUCT_NOT_FOUND",
				err.Error(),
				nil,
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				err.Error(),
				nil,
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		toOrderResponse(order),
	)
}

// Get godoc
// @Summary      Get an order by ID
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Security     BearerAuth
// @Success      200  {object}  OrderResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /orders/{id} [get]
func (h *OrderHandler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
			nil,
		)
		return
	}

	order, err := h.orderService.Get(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"ORDER_NOT_FOUND",
				err.Error(),
				nil,
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
			nil,
		)
		return
	}

	if order.UserID != userID {
		writeError(
			w,
			http.StatusNotFound,
			"ORDER_NOT_FOUND",
			"order not found",
			nil,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		toOrderResponse(order),
	)
}

// ListByUser godoc
// @Summary      List authenticated user's orders
// @Tags         orders
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array}   OrderResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /orders [get]
func (h *OrderHandler) ListByUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
			nil,
		)
		return
	}

	orders, err := h.orderService.ListByUser(
		r.Context(),
		userID,
	)

	if err != nil {
		log.Printf("ListByUser error: %v", err)

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
			nil,
		)
		return
	}

	out := make([]OrderResponse, 0, len(orders))

	for _, order := range orders {
		out = append(
			out,
			toOrderResponse(order),
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		out,
	)
}

// Cancel godoc
// @Summary      Cancel an order
// @Tags         orders
// @Produce      json
// @Param        id   path      string  true  "Order ID"
// @Security     BearerAuth
// @Success      200  {object}  OrderResponse
// @Failure      401  {object}  ErrorResponse
// @Failure      404  {object}  ErrorResponse
// @Failure      500  {object}  ErrorResponse
// @Router       /orders/{id}/cancel [post]
func (h *OrderHandler) Cancel(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
			nil,
		)
		return
	}

	order, err := h.orderService.Get(
		r.Context(),
		id,
	)
	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"ORDER_NOT_FOUND",
				err.Error(),
				nil,
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
			nil,
		)
		return
	}

	if order.UserID != userID {
		writeError(
			w,
			http.StatusNotFound,
			"ORDER_NOT_FOUND",
			"order not found",
			nil,
		)
		return
	}

	order, err = h.orderService.Cancel(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, domain.ErrOrderNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"ORDER_NOT_FOUND",
				err.Error(),
				nil,
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
			nil,
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		toOrderResponse(order),
	)
}
