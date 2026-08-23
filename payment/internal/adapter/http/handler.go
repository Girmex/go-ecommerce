package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/adapter/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/payment/internal/port"
)

type PaymentHandler struct {
	paymentService port.PaymentService
	validator      *validator.Validate
}

func NewPaymentHandler(
	paymentService port.PaymentService,
) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
		validator:      validator.New(),
	}
}

// Charge godoc
// @Summary      Charge a payment for an order
// @Description  Creates a payment and attempts to charge it through the payment gateway.
// @Tags         payments
// @Accept       json
// @Produce      json
// @Param        payment  body      ChargeRequest  true  "Charge payload"
// @Security     BearerAuth
// @Success      201      {object}  PaymentResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      401      {object}  ErrorResponse
// @Failure      402      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /payments [post]
func (h *PaymentHandler) Charge(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
		)
		return
	}

	var req ChargeRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"INVALID_BODY",
			"request body is not valid JSON",
		)
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeError(
			w,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"request failed validation",
		)
		return
	}

	payment, err := h.paymentService.Charge(
		r.Context(),
		port.ChargeInput{
			OrderID: req.OrderID,
			UserID:  userID,
			Amount:  req.Amount,
			Method:  req.Method,
		},
	)

	if err != nil {
		if errors.Is(err, domain.ErrPaymentDeclined) {
			writeError(
				w,
				http.StatusPaymentRequired,
				"PAYMENT_DECLINED",
				"payment was declined",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		toPaymentResponse(payment),
	)
}

// Get godoc
// @Summary      Get a payment by ID
// @Tags         payments
// @Produce      json
// @Param        id  path  string  true  "Payment ID"
// @Security     BearerAuth
// @Success      200  {object} PaymentResponse
// @Failure      401  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments/{id} [get]
func (h *PaymentHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
		)
		return
	}

	id := chi.URLParam(r, "id")

	payment, err := h.paymentService.Get(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
		)
		return
	}

	if payment.UserID != userID {
		writeError(
			w,
			http.StatusNotFound,
			"PAYMENT_NOT_FOUND",
			"payment not found",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		toPaymentResponse(payment),
	)
}

// List godoc
// @Summary      List authenticated user's payments
// @Tags         payments
// @Produce      json
// @Security     BearerAuth
// @Success      200  {array} PaymentResponse
// @Failure      401  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments [get]
func (h *PaymentHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
		)
		return
	}

	payments, err := h.paymentService.List(
		r.Context(),
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			err.Error(),
		)
		return
	}

	out := make([]PaymentResponse, 0)

	for _, payment := range payments {
		if payment.UserID != userID {
			continue
		}

		out = append(
			out,
			toPaymentResponse(payment),
		)
	}

	writeJSON(
		w,
		http.StatusOK,
		out,
	)
}

// Refund godoc
// @Summary      Refund a payment
// @Tags         payments
// @Produce      json
// @Param        id  path  string  true  "Payment ID"
// @Security     BearerAuth
// @Success      200  {object} PaymentResponse
// @Failure      401  {object} ErrorResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments/{id}/refund [post]
func (h *PaymentHandler) Refund(
	w http.ResponseWriter,
	r *http.Request,
) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(
			w,
			http.StatusUnauthorized,
			"UNAUTHORIZED",
			"authenticated user not found",
		)
		return
	}

	id := chi.URLParam(r, "id")

	payment, err := h.paymentService.Get(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
		)
		return
	}

	if payment.UserID != userID {
		writeError(
			w,
			http.StatusNotFound,
			"PAYMENT_NOT_FOUND",
			"payment not found",
		)
		return
	}

	payment, err = h.paymentService.Refund(
		r.Context(),
		id,
	)

	if err != nil {
		if errors.Is(err, domain.ErrPaymentNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PAYMENT_NOT_FOUND",
				"payment not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		toPaymentResponse(payment),
	)
}

func writeJSON(
	w http.ResponseWriter,
	status int,
	data any,
) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(data)
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
) {
	writeJSON(
		w,
		status,
		ErrorResponse{
			Code:    code,
			Message: message,
		},
	)
}
