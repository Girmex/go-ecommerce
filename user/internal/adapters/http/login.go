package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http/dto"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
)

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid login data")
		return
	}

	result, err := h.userService.Login(
		r.Context(),
		req.Email,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrInvalidCredentials):
			writeError(w, http.StatusUnauthorized, "invalid email or password")

		case errors.Is(err, application.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid login data")

		default:
			writeError(w, http.StatusInternalServerError, "failed to login")
		}

		return
	}

	response := dto.LoginResponse{
		ID:    result.User.ID,
		Name:  result.User.Name,
		Email: result.User.Email,
		Token: result.Token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}