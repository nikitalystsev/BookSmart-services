package intfRepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

//go:generate mockgen -source=IRatingRepo.go -destination=../../../internal/tests/unitTests/serviceTests/mocks/mockRatingRepo.go --package=mocks

type IRatingRepo interface {
	Create(ctx context.Context, rating *models.RatingModel) error
	GetByReaderAndBook(ctx context.Context, readerID uuid.UUID, bookID uuid.UUID) (*models.RatingModel, error)
	GetByBookID(ctx context.Context, bookID uuid.UUID, limit, offset int) ([]*models.RatingModel, error)
}
