package models

// CreateSessionResponse represents the response for session creation
//
// swagger:model CreateSessionResponse
type CreateSessionResponse struct {
	SessionID string `json:"session_id"`
}

// GetVotesResponse represents the response containing votes for a session
//
// swagger:model GetVotesResponse
type GetVotesResponse struct {
	Votes []Vote `json:"votes"`
}

// GetAggregatedScoresResponse represents the response containing aggregated scores
//
// swagger:model GetAggregatedScoresResponse
type GetAggregatedScoresResponse struct {
	Scores []ProductScore `json:"scores"`
}

// ErrorResponse represents a standard error response
//
// swagger:model ErrorResponse
type ErrorResponse struct {
	Message string `json:"message"`
}

// EmptyResponse represents an empty response
//
// swagger:model EmptyResponse
type EmptyResponse struct{}
