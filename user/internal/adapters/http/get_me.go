package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http/dto"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"
)
// @Summary Get current user
// @Description Returns the currently authenticated user.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/me [get]
func (h *Handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserID(r.Context())
	if !ok {
		writeError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	user, err := h.userService.GetUser(r.Context(), userID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			writeError(w, http.StatusNotFound, "user not found")
			return
		}

		writeError(w, http.StatusInternalServerError, "failed to get user")
		return
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
