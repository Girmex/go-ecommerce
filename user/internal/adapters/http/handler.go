package http

import (
	"encoding/json"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
)

type Handler struct {
	userService *application.UserService
}

func NewHandler(userService *application.UserService) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}
