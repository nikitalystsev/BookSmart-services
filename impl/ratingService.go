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

func (rs *RatingService) Create(ctx context.Context, rating *models.RatingModel) error {
	rs.logger.Info("attempting to create rating")

	existingRating, err := rs.ratingRepo.GetByReaderAndBook(ctx, rating.ReaderID, rating.BookID)
	if err != nil && !errors.Is(err, errs.ErrRatingDoesNotExists) {
		rs.logger.Errorf("error getting existing rating: %v", err)
		return err
	}

	if existingRating != nil {
		rs.logger.Warn("rating already exists")
		return errs.ErrRatingAlreadyExist
	}

	existingReservation, err := rs.reservationRepo.GetByReaderAndBook(ctx, rating.ReaderID, rating.BookID)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error getting existing reservation: %v", err)
		return err
	}

	if existingReservation == nil {
		rs.logger.Warn("reservation not found")
		return errs.ErrReservationDoesNotExists
	}

	if err = rs.ratingRepo.Create(ctx, rating); err != nil {
		rs.logger.Errorf("error creating rating: %v", err)
		return err
	}

	rs.logger.Infof("succesfully created rating")

	return nil
}

func (rs *RatingService) GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.RatingModel, error) {
	rs.logger.Infof("attempting to get ratings with bookID: %s", bookID)

	ratings, err := rs.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrRatingDoesNotExists) {
		rs.logger.Errorf("error getting ratings: %v", err)
		return nil, err
	}

	if errors.Is(err, errs.ErrRatingDoesNotExists) || len(ratings) == 0 {
		rs.logger.Warn("ratings not found")
		return nil, errs.ErrRatingDoesNotExists
	}

	rs.logger.Infof("succesfully getting ratings by bookID: %s", bookID)

	return ratings, nil
}

func (rs *RatingService) GetAvgRatingByBookID(ctx context.Context, bookID uuid.UUID) (float32, error) {
	rs.logger.Infof("attempting to get avg rating with bookID: %s", bookID)

	ratings, err := rs.ratingRepo.GetByBookID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrRatingDoesNotExists) {
		rs.logger.Errorf("error getting ratings: %v", err)
		return -1, err
	}

	if errors.Is(err, errs.ErrRatingDoesNotExists) || len(ratings) == 0 {
		rs.logger.Warn("ratings not found")
		return -1, errs.ErrRatingDoesNotExists
	}

	var total float32
	for _, rating := range ratings {
		total += float32(rating.Rating)
	}

	avgRating := total / float32(len(ratings))

	rs.logger.Infof("average rating for bookID %s: %f", bookID, avgRating)

	return avgRating, nil
}
