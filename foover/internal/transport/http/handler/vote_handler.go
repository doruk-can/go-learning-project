package handler

import (
	"encoding/json"
	"foover/internal/models"
	"foover/internal/service"
	"github.com/gorilla/mux"
	"log"
	"log/slog"
	"net/http"
)

// SaveVoteHandler handles saving or updating a vote
// @Summary Save or update a vote
// @Description Stores or updates a product vote for a given session ID.
// @Tags votes
// @Accept json
// @Produce json
// @Param vote body models.SaveVoteRequest true "Vote object that needs to be added or updated"
// @Success 201 {object} models.EmptyResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /votes [post]
func SaveVoteHandler(voteService service.VoteService, productService service.ProductService, sessionService service.SessionService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var voteReq models.SaveVoteRequest
		if err := json.NewDecoder(r.Body).Decode(&voteReq); err != nil {
			logger.Error("Invalid request payload", "error", err)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid request payload")
			return
		}

		// Validate the request
		if err := ValidateStruct(voteReq); err != nil {
			logger.Error("Validation failed", "error", err)
			writeErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		ctx := r.Context()
		// Validate session ID
		sessionExists, err := sessionService.SessionExists(ctx, voteReq.SessionID)
		if err != nil {
			logger.Error("Error validating session ID", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to validate session ID")
			return
		}
		if !sessionExists {
			logger.Warn("Invalid session ID", "sessionID", voteReq.SessionID)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid session ID")
			return
		}

		// Validate product ID
		isValidProduct, err := productService.IsValidProductID(ctx, voteReq.ProductID)
		if err != nil {
			logger.Error("Error validating product ID", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to validate product ID")
			return
		}
		if !isValidProduct {
			logger.Warn("Invalid product ID", "productID", voteReq.ProductID)
			writeErrorResponse(w, http.StatusBadRequest, "Invalid product ID")
			return
		}

		vote := models.Vote{
			SessionID: voteReq.SessionID,
			ProductID: voteReq.ProductID,
			Score:     voteReq.Score,
		}

		if err := voteService.SaveVote(ctx, vote); err != nil {
			logger.Error("Failed to save vote", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to save vote")
			return
		}

		w.WriteHeader(http.StatusCreated)
		logger.Info("Successfully saved vote", "sessionID", voteReq.SessionID, "productID", voteReq.ProductID)
	}
}

// GetVotesHandler retrieves votes for a given session ID
// @Summary Get votes by session ID
// @Description Retrieves existing votes for products for a given session ID.
// @Tags votes
// @Produce json
// @Param session_id path string true "The session ID"
// @Success 200 {object} models.GetVotesResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /votes/{session_id} [get]
func GetVotesHandler(voteService service.VoteService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sessionID := vars["session_id"]

		logger.Info("Received request to get votes", "sessionID", sessionID)
		ctx := r.Context()
		votes, err := voteService.GetVotesBySessionID(ctx, sessionID)
		if err != nil {
			log.Printf("Error getting votes: %v", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get votes")
			return
		}

		// Returning empty array if no votes found
		if votes == nil {
			votes = []models.Vote{}
		}

		response := models.GetVotesResponse{
			Votes: votes,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		logger.Info("Successfully retrieved and sent votes", "sessionID", sessionID)
	}
}
