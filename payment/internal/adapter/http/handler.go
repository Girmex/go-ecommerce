package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

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
// @Success      201      {object}  PaymentResponse
// @Failure      400      {object}  ErrorResponse
// @Failure      402      {object}  ErrorResponse
// @Failure      500      {object}  ErrorResponse
// @Router       /payments [post]
func (h *PaymentHandler) Charge(
	w http.ResponseWriter,
	r *http.Request,
) {
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
			UserID:  req.UserID,
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
// @Success      200  {object} PaymentResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments/{id} [get]
func (h *PaymentHandler) Get(
	w http.ResponseWriter,
	r *http.Request,
) {
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

	writeJSON(
		w,
		http.StatusOK,
		toPaymentResponse(payment),
	)
}

// List godoc
// @Summary      List payments
// @Tags         payments
// @Produce      json
// @Success      200  {array} PaymentResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments [get]
func (h *PaymentHandler) List(
	w http.ResponseWriter,
	r *http.Request,
) {
	payments, err := h.paymentService.List(
		r.Context(),
	)

	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"an internal server error occurred",
		)
		return
	}

	out := make([]PaymentResponse, 0, len(payments))

	for _, payment := range payments {
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
// @Success      200  {object} PaymentResponse
// @Failure      404  {object} ErrorResponse
// @Failure      500  {object} ErrorResponse
// @Router       /payments/{id}/refund [post]
func (h *PaymentHandler) Refund(
	w http.ResponseWriter,
	r *http.Request,
) {
	id := chi.URLParam(r, "id")

	payment, err := h.paymentService.Refund(
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
