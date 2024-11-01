package handler

import (
	"encoding/json"
	"foover/internal/models"
	"foover/internal/service"
	"io"
	"log/slog"
	"net/http"
)

// CreateSessionHandler handles session creation
// @Summary Create a new session
// @Description Generates a unique session ID.
// @Tags sessions
// @Accept json
// @Produce json
// @Param session body models.CreateSessionRequest false "Session creation request"
// @Success 200 {object} models.CreateSessionResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /sessions [post]
func CreateSessionHandler(sessionService service.SessionService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req models.CreateSessionRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil && err != io.EOF {
			logger.Error("Invalid request payload", "error", err)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		logger.Info("Creating session")
		ctx := r.Context()
		sessionID, err := sessionService.CreateSession(ctx)
		if err != nil {
			logger.Error("Failed to create session", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to create session")
			return
		}

		response := models.CreateSessionResponse{SessionID: sessionID}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		logger.Info("Successfully created session", "sessionID", sessionID)
	}
}
