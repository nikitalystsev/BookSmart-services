package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	"github.com/sirupsen/logrus"
)

type RatingService struct {
	ratingRepo      intfRepo.IRatingRepo
	reservationRepo intfRepo.IReservationRepo
	logger          *logrus.Entry
}

func NewRatingService(
	ratingRepo intfRepo.IRatingRepo,
	reservationRepo intfRepo.IReservationRepo,
	logger *logrus.Entry,
) intf.IRatingService {
	return &RatingService{
		ratingRepo:      ratingRepo,
		reservationRepo: reservationRepo,
		logger:          logger,
	}
}

// Create TODO логировать
func (rs *RatingService) Create(ctx context.Context, rating *models.RatingModel) error {
	existingRating, err := rs.ratingRepo.GetByReaderAndBook(ctx, rating.ReaderID, rating.BookID)
	if err != nil && errors.Is(err, errs.ErrRatingAlreadyExist) {
		return err
	}

	if existingRating != nil {
		return errs.ErrRatingAlreadyExist
	}

	existingReservation, err := rs.reservationRepo.GetByReaderAndBook(ctx, rating.ReaderID, rating.BookID)
	if err != nil && errors.Is(err, errs.ErrReservationDoesNotExists) {
		return err
	}

	if existingReservation == nil {
		return errs.ErrReservationDoesNotExists
	}

	if err = rs.ratingRepo.Create(ctx, rating); err != nil {
		return err
	}

	return nil
}

// GetByBookID TODO логировать
func (rs *RatingService) GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.RatingModel, error) {
	ratings, err := rs.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrRatingDoesNotExists) {
		return nil, err
	}

	if errors.Is(err, errs.ErrRatingDoesNotExists) || len(ratings) == 0 {
		return nil, errs.ErrRatingDoesNotExists
	}

	return ratings, nil
}
