package service

import (
	"context"
	"foover/internal/models"
	"foover/internal/store/mongo"
)

// voteService implements the VoteService interface
type voteService struct {
	store mongo.Store
}

type VoteService interface {
	SaveVote(ctx context.Context, vote models.Vote) error
	GetVotesBySessionID(ctx context.Context, sessionID string) ([]models.Vote, error)
}

// NewVoteService creates a new VoteService
func NewVoteService(store mongo.Store) VoteService {
	return &voteService{
		store: store,
	}
}

// SaveVote stores or updates a vote for a given session ID and product ID
func (v *voteService) SaveVote(ctx context.Context, vote models.Vote) error {
	return v.store.SaveVote(ctx, vote)
}

// GetVotesBySessionID retrieves all votes associated with a session ID
func (v *voteService) GetVotesBySessionID(ctx context.Context, sessionID string) ([]models.Vote, error) {
	return v.store.GetVotesBySessionID(ctx, sessionID)
}
