package impl

import (
	"context"
	"errors"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
)

type RatingService struct {
	ratingRepo      intfRepo.IRatingRepo
	reservationRepo intfRepo.IReservationRepo
}

func NewRatingService(ratingRepo intfRepo.IRatingRepo, reservationRepo intfRepo.IReservationRepo) intf.IRatingService {
	return &RatingService{ratingRepo: ratingRepo, reservationRepo: reservationRepo}
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
