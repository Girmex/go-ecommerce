package handlers

import (
	"net/http"

	authproto "github.com/Girmex/go-ecommerce/microservices/auth/proto"
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
	w.Write([]byte("register"))
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("login"))
}