package service

import (
	"context"
	"foover/internal/models"
	"foover/internal/store/mongo"
)

// aggregationService implements the AggregationService interface
type aggregationService struct {
	store mongo.Store
}

type AggregationService interface {
	GetAggregatedProductScores(ctx context.Context) ([]models.ProductScore, error)
}

// NewAggregationService creates a new AggregationService
func NewAggregationService(store mongo.Store) AggregationService {
	return &aggregationService{
		store: store,
	}
}

// GetAggregatedProductScores retrieves aggregated average scores for all products
func (a *aggregationService) GetAggregatedProductScores(ctx context.Context) ([]models.ProductScore, error) {
	return a.store.GetAggregatedProductScores(ctx)
}
