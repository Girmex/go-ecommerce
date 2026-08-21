package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapters/http/dto"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/go-playground/validator/v10"
)

type Handler struct {
	userService *application.UserService
	validator   *validator.Validate
}

func NewHandler(userService *application.UserService) *Handler {
	return &Handler{
		userService: userService,
		validator:   validator.New(),
	}
}

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	_ = json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
	})
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		validationErr := err.(validator.ValidationErrors)[0]

		switch validationErr.Field() {
		case "Name":
			if validationErr.Tag() == "required" {
				writeError(w, http.StatusBadRequest, "name is required")
			} else {
				writeError(w, http.StatusBadRequest, "name must be at least 2 characters")
			}

		case "Email":
			if validationErr.Tag() == "required" {
				writeError(w, http.StatusBadRequest, "email is required")
			} else {
				writeError(w, http.StatusBadRequest, "email must be a valid email address")
			}

		case "Password":
			if validationErr.Tag() == "required" {
				writeError(w, http.StatusBadRequest, "password is required")
			} else {
				writeError(w, http.StatusBadRequest, "password must be at least 8 characters")
			}

		default:
			writeError(w, http.StatusBadRequest, "invalid user data")
		}

		return
	}

	user, err := h.userService.CreateUser(
		r.Context(),
		req.Name,
		req.Email,
		req.Password,
	)
	if err != nil {
		switch {
		case errors.Is(err, application.ErrUserAlreadyExists):
			writeError(w, http.StatusConflict, "user already exists")

		case errors.Is(err, application.ErrInvalidInput):
			writeError(w, http.StatusBadRequest, "invalid user input")

		default:
			writeError(w, http.StatusInternalServerError, "failed to create user")
		}

		return
	}

	response := dto.UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)
}
