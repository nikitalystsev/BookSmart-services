package intfRepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

//go:generate mockgen -source=ILibCardRepo.go -destination=../../../internal/tests/unitTests/serviceTests/mocks/mockLibCardRepo.go --package=mocks

type ILibCardRepo interface {
	Create(ctx context.Context, libCard *models.LibCardModel) error
	GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error)
	GetByNum(ctx context.Context, libCardNum string) (*models.LibCardModel, error)
	Update(ctx context.Context, libCard *models.LibCardModel) error
}
