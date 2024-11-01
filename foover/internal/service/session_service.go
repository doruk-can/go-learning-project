package service

import (
	"context"
	"foover/internal/store/mongo"
)

type SessionService interface {
	CreateSession(ctx context.Context) (string, error)
	SessionExists(ctx context.Context, sessionID string) (bool, error)
}

// sessionService implements the SessionService interface
type sessionService struct {
	store mongo.Store
}

// NewSessionService creates a new SessionService
func NewSessionService(store mongo.Store) SessionService {
	return &sessionService{
		store: store,
	}
}

// CreateSession creates a new session
func (s *sessionService) CreateSession(ctx context.Context) (string, error) {
	return s.store.CreateSession(ctx)
}

// SessionExists checks if a session with the given sessionID exists
func (s *sessionService) SessionExists(ctx context.Context, sessionID string) (bool, error) {
	return s.store.SessionExists(ctx, sessionID)
}
