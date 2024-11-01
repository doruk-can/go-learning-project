package models

// SaveVoteRequest represents the request payload for saving a vote
//
// swagger:model SaveVoteRequest
type SaveVoteRequest struct {
	SessionID string `json:"session_id" validate:"required,uuid4"`
	ProductID string `json:"product_id" validate:"required,uuid4"`
	Score     int    `json:"score" validate:"required,min=1,max=5"`
}

// CreateSessionRequest represents the request to create a session
//
// swagger:model CreateSessionRequest
type CreateSessionRequest struct { // Created for extension purposes

}
