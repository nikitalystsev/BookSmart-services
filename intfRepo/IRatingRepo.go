package intfRepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type IRatingRepo interface {
	GetByReaderAndBook(ctx context.Context, readerID uuid.UUID, bookID uuid.UUID) (*models.RatingModel, error)
	Create(ctx context.Context, rating *models.RatingModel) error
}
