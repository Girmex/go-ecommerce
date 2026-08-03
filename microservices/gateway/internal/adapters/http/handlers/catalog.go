package handlers

import (
	"encoding/json"
	"net/http"

	catalogproto "github.com/Girmex/go-ecommerce/microservices/catalog/proto"
	cataloggrpc "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/grpc/catalog"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/middleware"
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

func (h *CatalogHandler) ListProducts(
	w http.ResponseWriter,
	r *http.Request,
) {

	resp, err := h.client.ListProducts(
		r.Context(),
		&emptypb.Empty{},
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(resp)
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
