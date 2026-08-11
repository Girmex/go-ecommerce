package handlers

import (
	"encoding/json"
	"net/http"

	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
	httpadapter "github.com/Girmex/go-ecommerce/microservices/gateway/internal/adapters/http"
	"github.com/Girmex/go-ecommerce/microservices/gateway/internal/dto"
)

type AuthHandler struct {
	AuthClient authproto.AuthServiceClient
}

func NewAuthHandler(client authproto.AuthServiceClient) *AuthHandler {
	return &AuthHandler{
		AuthClient: client,
	}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {

	var req dto.RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.AuthClient.Register(
		r.Context(),
		&authproto.RegisterRequest{
			Name:     req.Name,
			Email:    req.Email,
			Password: req.Password,
			Phone:    req.Phone,
		},
	)

	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {

	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp, err := h.AuthClient.Login(
		r.Context(),
		&authproto.LoginRequest{
			Email:    req.Email,
			Password: req.Password,
		},
	)

	if err != nil {
		httpadapter.WriteGRPCError(w, err)
		return
	}

	httpadapter.WriteJSON(w, http.StatusOK, resp)
}
