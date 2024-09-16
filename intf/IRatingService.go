package intf

import (
	"context"
	"github.com/nikitalystsev/BookSmart-services/core/models"
)

type IRatingService interface {
	Create(ctx context.Context, rating *models.RatingModel) error
}
