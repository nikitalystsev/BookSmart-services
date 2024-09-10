package intf

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type IReservationService interface {
	Create(ctx context.Context, readerID, bookID uuid.UUID) error
	Update(ctx context.Context, reservation *models.ReservationModel) error
	GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.ReservationModel, error)
	GetAllReservationsByReaderID(ctx context.Context, readerID uuid.UUID) ([]*models.ReservationModel, error)
	GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error)
}
