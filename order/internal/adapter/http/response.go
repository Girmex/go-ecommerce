package http

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details any         `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		return
	}
}

func writeError(
	w http.ResponseWriter,
	status int,
	code string,
	message string,
	details any,
) {
	writeJSON(w, status, ErrorResponse{
		Code:    code,
		Message: message,
		Details: details,
	})
}