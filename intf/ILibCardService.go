package intf

import (
	"context"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type ILibCardService interface {
	Create(ctx context.Context, readerID uuid.UUID) error
	Update(ctx context.Context, libCard *models.LibCardModel) error
	GetByReaderID(ctx context.Context, readerID uuid.UUID) (*models.LibCardModel, error)
}
