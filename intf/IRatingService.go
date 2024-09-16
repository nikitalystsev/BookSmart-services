package intf

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type IRatingService interface {
	Create(ctx context.Context, rating *models.RatingModel) error
	GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.RatingModel, error)
	GetAvgRatingByBookID(ctx context.Context, bookID uuid.UUID) (float32, error)
}
