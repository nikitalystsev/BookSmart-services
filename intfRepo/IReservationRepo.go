package intfRepo

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

//go:generate mockgen -source=IReservationRepo.go -destination=../../../internal/tests/unitTests/serviceTests/mocks/mockReservationRepo.go --package=mocks

type IReservationRepo interface {
	Create(ctx context.Context, reservation *models.ReservationModel) error
	GetByReaderAndBook(ctx context.Context, readerID, bookID uuid.UUID) ([]*models.ReservationModel, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error)
	GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.ReservationModel, error)
	Update(ctx context.Context, reservation *models.ReservationModel) error
	GetExpiredByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
	GetActiveByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
	GetByReaderID(ctx context.Context, readerID uuid.UUID, limit, offset int) ([]*models.ReservationModel, error)
}
