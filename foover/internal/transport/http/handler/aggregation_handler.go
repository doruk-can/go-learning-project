package handler

import (
	"encoding/json"
	"foover/internal/models"
	"foover/internal/service"
	"log/slog"
	"net/http"
)

// GetAggregatedScoresHandler retrieves aggregated product scores
// @Summary Get aggregated product scores
// @Description Retrieves aggregated average scores for products across all session IDs.
// @Tags aggregation
// @Produce json
// @Success 200 {object} models.GetAggregatedScoresResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /aggregated-scores [get]
func GetAggregatedScoresHandler(aggregationService service.AggregationService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		scores, err := aggregationService.GetAggregatedProductScores(ctx)
		if err != nil {
			logger.Error("Failed to get aggregated scores", "error", err)
			writeErrorResponse(w, http.StatusInternalServerError, "Failed to get aggregated scores")
			return
		}

		// Return an empty array if no scores are found
		if scores == nil {
			scores = []models.ProductScore{}
		}

		response := models.GetAggregatedScoresResponse{
			Scores: scores,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		logger.Info("Successfully retrieved and sent aggregated scores")
	}
}
