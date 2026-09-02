package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/adapter/http/middleware"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/application"
	"github.com/Girmex/go-ecommerce-app/chi-microservice/user/internal/domain"

	"github.com/go-playground/validator/v10"
)

type Handler struct {
	userService *application.UserService
	validator   *validator.Validate
}

type ErrorResponse struct {
	Message          string       `json:"message"`
	ValidationErrors []FieldError `json:"validation_errors,omitempty"`
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

func writeError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	_ = json.NewEncoder(w).Encode(ErrorResponse{
		Message: message,
	})
}

// CreateUser godoc
// @Summary Create user
// @Description Creates a new user.
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserRequest true "User"
// @Success 201 {object} UserResponse
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users [post]
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeValidationError(w, err)
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

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	_ = json.NewEncoder(w).Encode(response)
}

// Login godoc
// @Summary Login
// @Description Authenticates a user and returns a JWT.
// @Tags Users
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /users/login [post]
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := h.validator.Struct(req); err != nil {
		writeValidationError(w, err)
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

	response := LoginResponse{
		ID:    result.User.ID,
		Name:  result.User.Name,
		Email: result.User.Email,
		Token: result.Token,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}

// @Summary Get current user
// @Description Returns the currently authenticated user.
// @Tags Users
// @Produce json
// @Security BearerAuth
// @Success 200 {object} UserResponse
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

	response := UserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	_ = json.NewEncoder(w).Encode(response)
}
