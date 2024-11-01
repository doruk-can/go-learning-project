package handler

import (
	"encoding/json"
	"foover/internal/models"
	"net/http"
)

func writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	errorResponse := models.ErrorResponse{
		Message: message,
	}
	json.NewEncoder(w).Encode(errorResponse)
}
