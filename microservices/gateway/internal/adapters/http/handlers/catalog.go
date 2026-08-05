package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	catalogproto "github.com/Girmex/go-ecommerce/microservices/catalog/proto"
	cataloggrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/catalog"
	httpadapter "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/http"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
	"github.com/go-chi/chi/v5"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CatalogHandler struct {
	client *cataloggrpc.Client
}

func NewCatalogHandler(client *cataloggrpc.Client) *CatalogHandler {
	return &CatalogHandler{
		client: client,
	}
}

func (h *CatalogHandler) CreateCategory(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims

	grpcReq := &catalogproto.CreateCategoryRequest{
		Name:         req.Name,
		ParentId:     req.ParentID,
		ImageUrl:     req.ImageURL,
		DisplayOrder: req.DisplayOrder,
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.CreateCategory(
		ctx,
		grpcReq,
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) ListCategories(
	w http.ResponseWriter,
	r *http.Request,
) {

	resp, err := h.client.ListCategories(
		r.Context(),
		&emptypb.Empty{},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) UpdateCategory(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateCategoryRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.UpdateCategory(
		ctx,
		&catalogproto.UpdateCategoryRequest{
			Id:           uint32(id),
			Name:         req.Name,
			ParentId:     req.ParentID,
			ImageUrl:     req.ImageURL,
			DisplayOrder: req.DisplayOrder,
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) GetCategory(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	resp, err := h.client.GetCategory(
		r.Context(),
		&catalogproto.GetCategoryRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) DeleteCategory(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid category id", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	_, err = h.client.DeleteCategory(
		ctx,
		&catalogproto.GetCategoryRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteNoContent(w)
}

func (h *CatalogHandler) ListProducts(
	w http.ResponseWriter,
	r *http.Request,
) {

	resp, err := h.client.ListProducts(
		r.Context(),
		&emptypb.Empty{},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) CreateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {
	var req dto.CreateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	grpcReq := &catalogproto.CreateProductRequest{
		Name:        req.Name,
		Description: req.Description,
		CategoryId:  req.CategoryID,
		ImageUrl:    req.ImageURL,
		Price:       req.Price,
		Stock:       req.Stock,
		UserId:      uint32(claims.UserID),
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.CreateProduct(
		ctx,
		grpcReq,
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) GetProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	resp, err := h.client.GetProduct(
		r.Context(),
		&catalogproto.GetProductRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) GetSellerProducts(
	w http.ResponseWriter,
	r *http.Request,
) {

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.GetSellerProducts(
		ctx,
		&catalogproto.GetSellerProductsRequest{
			UserId: uint32(claims.UserID),
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *CatalogHandler) UpdateProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims // you'll use this later when ownership checks are added

	grpcReq := &catalogproto.UpdateProductRequest{
		Id: uint32(id),
	}

	if req.Name != nil {
		grpcReq.Name = req.Name
	}

	if req.Description != nil {
		grpcReq.Description = req.Description
	}

	if req.CategoryID != nil {
		grpcReq.CategoryId = req.CategoryID
	}

	if req.ImageURL != nil {
		grpcReq.ImageUrl = req.ImageURL
	}

	if req.Price != nil {
		grpcReq.Price = req.Price
	}

	if req.Stock != nil {
		grpcReq.Stock = req.Stock
	}

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.UpdateProduct(
		ctx,
		grpcReq,
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}
func (h *CatalogHandler) DeleteProduct(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	_, err = h.client.DeleteProduct(
		ctx,
		&catalogproto.GetProductRequest{
			Id: uint32(id),
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteNoContent(w)
}

func (h *CatalogHandler) UpdateProductStock(
	w http.ResponseWriter,
	r *http.Request,
) {

	idParam := chi.URLParam(r, "id")

	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		http.Error(w, "invalid product id", http.StatusBadRequest)
		return
	}

	var req dto.UpdateStockRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	claims, ok := middleware.Claims(r.Context())
	if !ok {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	_ = claims

	authHeader := r.Header.Get("Authorization")

	ctx := metadata.AppendToOutgoingContext(
		r.Context(),
		"authorization",
		authHeader,
	)

	resp, err := h.client.UpdateProductStock(
		ctx,
		&catalogproto.UpdateStockRequest{
			Id:    uint32(id),
			Stock: req.Stock,
		},
	)
	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}
