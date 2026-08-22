package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/adapters/http/dto"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/domain"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/product/internal/ports"
)

type Handler struct {
	productService ports.ProductService
	validator      *validator.Validate
}

func NewHandler(productService ports.ProductService) *Handler {
	return &Handler{
		productService: productService,
		validator:      validator.New(),
	}
}

// Create godoc
// @Summary      Create a product
// @Description  Creates a new product.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        product  body      dto.CreateProductRequest  true  "Product payload"
// @Success      201      {object}  dto.ProductResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /products [post]
func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateProductRequest

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

	product, err := h.productService.Create(
		r.Context(),
		ports.CreateProductInput{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
		},
	)
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	writeJSON(
		w,
		http.StatusCreated,
		dto.ToProductResponse(product),
	)
}

// Get godoc
// @Summary      Get a product by ID
// @Description  Returns a product by its ID.
// @Tags         products
// @Produce      json
// @Param        id  path      string  true  "Product ID"
// @Success      200  {object}  dto.ProductResponse
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /products/{id} [get]
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	product, err := h.productService.Get(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToProductResponse(product),
	)
}

// List godoc
// @Summary      List products
// @Description  Returns all products.
// @Tags         products
// @Produce      json
// @Success      200  {array}   dto.ProductResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /products [get]
func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	products, err := h.productService.List(r.Context())
	if err != nil {
		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	response := make([]dto.ProductResponse, 0, len(products))

	for _, product := range products {
		response = append(
			response,
			dto.ToProductResponse(product),
		)
	}

	writeJSON(w, http.StatusOK, response)
}

// Update godoc
// @Summary      Update a product
// @Description  Updates an existing product.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                  true  "Product ID"
// @Param        product  body      dto.UpdateProductRequest true  "Product fields"
// @Success      200      {object}  dto.ProductResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /products/{id} [put]
func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateProductRequest

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

	product, err := h.productService.Update(
		r.Context(),
		id,
		ports.UpdateProductInput{
			Name:        req.Name,
			Description: req.Description,
			Price:       req.Price,
			Stock:       req.Stock,
		},
	)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToProductResponse(product),
	)
}

// Delete godoc
// @Summary      Delete a product
// @Description  Deletes a product by its ID.
// @Tags         products
// @Produce      json
// @Param        id  path  string  true  "Product ID"
// @Success      204  "No Content"
// @Failure      404  {object}  dto.ErrorResponse
// @Failure      500  {object}  dto.ErrorResponse
// @Router       /products/{id} [delete]
func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	err := h.productService.Delete(r.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			writeError(
				w,
				http.StatusNotFound,
				"PRODUCT_NOT_FOUND",
				"product not found",
			)
			return
		}

		writeError(
			w,
			http.StatusInternalServerError,
			"INTERNAL_ERROR",
			"internal server error",
		)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ReserveStock godoc
// @Summary      Reserve product stock
// @Description  Reserves product stock for an order.
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                   true  "Product ID"
// @Param        reserve  body      dto.ReserveStockRequest true  "Stock reservation"
// @Success      200      {object}  dto.ProductResponse
// @Failure      400      {object}  dto.ErrorResponse
// @Failure      404      {object}  dto.ErrorResponse
// @Failure      409      {object}  dto.ErrorResponse
// @Failure      500      {object}  dto.ErrorResponse
// @Router       /products/{id}/reserve [post]
func (h *Handler) ReserveStock(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.ReserveStockRequest

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

	product, err := h.productService.ReserveStock(
		r.Context(),
		id,
		req.Quantity,
	)
	if err != nil {
		switch {
		case errors.Is(err, domain.ErrProductNotFound):
			writeError(
				w,
				http.StatusNotFound,
				"PRODUCT_NOT_FOUND",
				"product not found",
			)

		case errors.Is(err, domain.ErrInsufficientStock):
			writeError(
				w,
				http.StatusConflict,
				"INSUFFICIENT_STOCK",
				"insufficient stock",
			)

		default:
			writeError(
				w,
				http.StatusInternalServerError,
				"INTERNAL_ERROR",
				"internal server error",
			)
		}

		return
	}

	writeJSON(
		w,
		http.StatusOK,
		dto.ToProductResponse(product),
	)
}

func writeJSON(w http.ResponseWriter, status int, data any) {
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
		dto.ErrorResponse{
			Code:    code,
			Message: message,
		},
	)
}