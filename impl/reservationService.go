package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	"github.com/nikitalystsev/BookSmart-services/pkg/transact"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	ReservationIssued   = "Issued"
	ReservationExtended = "Extended"
	ReservationExpired  = "Expired"
	ReservationClosed   = "Closed"
)

const (
	ReservationIssuePeriodDays     = 14
	ReservationExtensionPeriodDays = 7

	ReservationsPageLimit = 10
)

type ReservationService struct {
	reservationRepo    intfRepo.IReservationRepo
	bookRepo           intfRepo.IBookRepo
	readerRepo         intfRepo.IReaderRepo
	libCardRepo        intfRepo.ILibCardRepo
	transactionManager transact.ITransactionManager
	logger             *logrus.Entry
}

func NewReservationService(
	reservationRepo intfRepo.IReservationRepo,
	bookRepo intfRepo.IBookRepo,
	readerRepo intfRepo.IReaderRepo,
	libCardRepo intfRepo.ILibCardRepo,
	transactionManager transact.ITransactionManager,
	logger *logrus.Entry,
) intf.IReservationService {
	return &ReservationService{
		reservationRepo:    reservationRepo,
		bookRepo:           bookRepo,
		readerRepo:         readerRepo,
		libCardRepo:        libCardRepo,
		transactionManager: transactionManager,
		logger:             logger,
	}
}

func (rs *ReservationService) Create(ctx context.Context, readerID, bookID uuid.UUID) error {
	rs.logger.Info("starting reservation creation process")

	existingReader, err := rs.checkReaderCanCreateReservation(ctx, readerID)
	if err != nil {
		return err
	}

	existingBook, err := rs.checkBookCanBeReserved(ctx, bookID)
	if err != nil {
		return err
	}

	if err = rs.checkAgeLimit(existingReader, existingBook); err != nil {
		return err
	}

	if err = rs.create(ctx, readerID, bookID); err != nil {
		rs.logger.Errorf("error creating reservation: %v", err)
		return err
	}

	rs.logger.Info("reservation creation successful")

	return nil
}

func (rs *ReservationService) Update(ctx context.Context, reservation *models.ReservationModel, extentionPeriodDays int) error {
	if reservation == nil {
		rs.logger.Warn("reservation object is nil")
		return errs.ErrReservationObjectIsNil
	}

	rs.logger.Info("attempting to update reservation")

	if err := rs.checkValidLibCard(ctx, reservation.ReaderID); err != nil {
		return err
	}

	if err := rs.checkNoExpiredBooks(ctx, reservation.ReaderID); err != nil {
		return err
	}

	if err := rs.checkReservationState(reservation.State); err != nil {
		return err
	}

	if err := rs.checkBookIsCommon(ctx, reservation.BookID); err != nil {
		return err
	}

	if err := rs.checkExtentionPeriodDays(extentionPeriodDays); err != nil {
		return err
	}

	reservation.ReturnDate = reservation.ReturnDate.AddDate(0, 0, extentionPeriodDays)
	reservation.State = ReservationExtended

	rs.logger.Info("update reservation in repository")

	if err := rs.reservationRepo.Update(ctx, reservation); err != nil {
		rs.logger.Errorf("error updating reservation: %v", err)
		return err
	}

	rs.logger.Info("reservation update successful")

	return nil
}

// GetByBookID TODO добавить в схемы (протестировано)
func (rs *ReservationService) GetByBookID(ctx context.Context, bookID uuid.UUID) ([]*models.ReservationModel, error) {
	rs.logger.Infof("attempting to get reservation with bookID: %s", bookID)

	reservations, err := rs.reservationRepo.GetByBookID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error checking reservation existence: %v", err)
		return nil, err
	}

	if reservations == nil {
		rs.logger.Warnf("reservations with this bookID does not exist %s", bookID)
		return nil, errs.ErrReservationDoesNotExists
	}

	rs.logger.Infof("successfully getting reservation by bookID: %s", bookID)

	return reservations, nil

}

func (rs *ReservationService) GetByReaderID(ctx context.Context, readerID uuid.UUID, limit, offset int) ([]*models.ReservationModel, error) {
	reservations, err := rs.reservationRepo.GetByReaderID(ctx, readerID, limit, offset)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error checking reservations: %v", err)
		return nil, err
	}

	rs.logger.Info("successfully get reservations")

	return reservations, nil
}

func (rs *ReservationService) GetByID(ctx context.Context, ID uuid.UUID) (*models.ReservationModel, error) {
	rs.logger.Infof("attempting to get reservation with ID: %s", ID)

	reservation, err := rs.reservationRepo.GetByID(ctx, ID)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error checking reservation existence: %v", err)
		return nil, err
	}

	if reservation == nil {
		rs.logger.Warn("reservation with this ID does not exist")
		return nil, errs.ErrReservationDoesNotExists
	}

	rs.logger.Infof("successfully getting reservation by ID: %s", ID)

	return reservation, nil
}

func (rs *ReservationService) create(ctx context.Context, readerID, bookID uuid.UUID) error {
	return rs.transactionManager.Do(ctx, func(ctx context.Context) error {
		existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
		if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
			rs.logger.Errorf("error checking book existence: %v", err)
			return err
		}

		if err = rs.checkBookCopiesNumber(existingBook); err != nil {
			return err
		}

		newReservation := &models.ReservationModel{
			ID:         uuid.New(),
			ReaderID:   readerID,
			BookID:     bookID,
			IssueDate:  time.Now(),
			ReturnDate: time.Now().AddDate(0, 0, ReservationIssuePeriodDays),
			State:      ReservationIssued,
		}

		rs.logger.Info("creating reservation in repository")

		if err = rs.reservationRepo.Create(ctx, newReservation); err != nil {
			rs.logger.Errorf("error creating reservation: %v", err)
			return err
		}

		existingBook.CopiesNumber -= 1

		rs.logger.Info("updating book copiesNumber in repository")

		if err = rs.bookRepo.Update(ctx, existingBook); err != nil {
			rs.logger.Errorf("error updating book: %v", err)
			return err
		}

		rs.logger.Info("successfully updated book copiesNumber")

		return nil
	})
}

func (rs *ReservationService) checkReaderCanCreateReservation(ctx context.Context, readerID uuid.UUID) (*models.ReaderModel, error) {
	existingReader, err := rs.checkReaderExists(ctx, readerID)
	if err != nil {
		return nil, err
	}

	if err = rs.checkNoExpiredBooks(ctx, readerID); err != nil {
		return nil, err
	}

	if err = rs.checkActiveReservationsLimit(ctx, readerID); err != nil {
		return nil, err
	}

	if err = rs.checkValidLibCard(ctx, readerID); err != nil {
		return nil, err
	}

	rs.logger.Info("reader is valid")

	return existingReader, nil
}

func (rs *ReservationService) checkBookCanBeReserved(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	existingBook, err := rs.checkBookExists(ctx, bookID)
	if err != nil {
		return nil, err
	}

	if err = rs.checkBookIsCommonOrRare(existingBook); err != nil {
		return nil, err
	}

	rs.logger.Info("book is valid")

	return existingBook, nil
}

func (rs *ReservationService) checkReaderExists(ctx context.Context, readerID uuid.UUID) (*models.ReaderModel, error) {
	existingReader, err := rs.readerRepo.GetByID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrReaderDoesNotExists) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return nil, err
	}
	if existingReader == nil {
		rs.logger.Warn("reader with this ID does not exist")
		return nil, errs.ErrReaderDoesNotExists
	}

	rs.logger.Info("reader exists")

	return existingReader, nil
}

func (rs *ReservationService) checkNoExpiredBooks(ctx context.Context, readerID uuid.UUID) error {
	expiredReservations, err := rs.reservationRepo.GetExpiredByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error checking expired book existence: %v", err)
		return err
	}

	if len(expiredReservations) > 0 {
		rs.logger.Warn("reader has expired books")
		return errs.ErrReaderHasExpiredBooks
	}

	rs.logger.Info("reader has not expired books")

	return nil
}

func (rs *ReservationService) checkActiveReservationsLimit(ctx context.Context, readerID uuid.UUID) error {
	activeReservations, err := rs.reservationRepo.GetActiveByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrReservationDoesNotExists) {
		rs.logger.Errorf("error checking active reservations: %v", err)
		return err
	}
	if len(activeReservations) >= MaxBooksPerReader {
		rs.logger.Warn("reader has reached the limit of active reservations")
		return errs.ErrReservationsLimitExceeded
	}

	rs.logger.Info("reader has not reached the limit of active reservations")

	return nil
}

func (rs *ReservationService) checkValidLibCard(ctx context.Context, readerID uuid.UUID) error {
	libCard, err := rs.libCardRepo.GetByReaderID(ctx, readerID)
	if err != nil && !errors.Is(err, errs.ErrLibCardDoesNotExists) {
		rs.logger.Errorf("error checking libCard existence: %v", err)
		return err
	}
	if libCard == nil {
		rs.logger.Warn("reader does not have libCard")
		return errs.ErrLibCardDoesNotExists
	}

	if !libCard.ActionStatus {
		rs.logger.Warn("reader has invalid libCard")
		return errs.ErrLibCardIsInvalid
	}

	rs.logger.Info("reader has valid libCard")

	return nil
}

func (rs *ReservationService) checkBookExists(ctx context.Context, bookID uuid.UUID) (*models.BookModel, error) {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return nil, err
	}
	if existingBook == nil {
		rs.logger.Warn("book with this ID does not exist")
		return nil, errs.ErrBookDoesNotExists
	}

	rs.logger.Info("book exists")

	return existingBook, nil
}

func (rs *ReservationService) checkBookCopiesNumber(book *models.BookModel) error {
	if book.CopiesNumber <= 0 {
		rs.logger.Warn("no copies of the book are available in the library")
		return errs.ErrBookNoCopiesNum
	}

	rs.logger.Info("book has copies available")

	return nil
}

func (rs *ReservationService) checkBookIsCommonOrRare(book *models.BookModel) error {
	if book.Rarity == BookRarityUnique {
		rs.logger.Warn("this book is unique and cannot be reserved")
		return errs.ErrUniqueBookNotReserved
	}

	rs.logger.Info("book is not unique")

	return nil
}

func (rs *ReservationService) checkAgeLimit(reader *models.ReaderModel, book *models.BookModel) error {
	if reader.Age < book.AgeLimit {
		rs.logger.Warn("reader does not meet the age requirement for this book")
		return errs.ErrReservationAgeLimit
	}

	rs.logger.Info("reader's age is appropriate")

	return nil
}

func (rs *ReservationService) checkBookIsCommon(ctx context.Context, bookID uuid.UUID) error {
	existingBook, err := rs.bookRepo.GetByID(ctx, bookID)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		rs.logger.Errorf("error checking book existence: %v", err)
		return err
	}

	if existingBook.Rarity == BookRarityRare || existingBook.Rarity == BookRarityUnique {
		rs.logger.Warn("rare and unique book cannot be renewed.")
		return errs.ErrRareAndUniqueBookNotExtended
	}

	rs.logger.Info("book's rarity is common")

	return nil
}

func (rs *ReservationService) checkReservationState(reservationState string) error {
	if reservationState == ReservationClosed {
		rs.logger.Warn("this reservation is already closed")
		return errs.ErrReservationIsAlreadyClosed
	}

	if reservationState == ReservationExpired {
		rs.logger.Warn("this reservation is already expired")
		return errs.ErrReservationIsAlreadyExpired
	}

	if reservationState == ReservationExtended {
		rs.logger.Warn("this reservation is already extended")
		return errs.ErrReservationIsAlreadyExtended
	}

	rs.logger.Info("reservation is only issued")

	return nil
}

func (rs *ReservationService) checkExtentionPeriodDays(extentionPeriodDays int) error {
	if extentionPeriodDays >= ReservationExtensionPeriodDays {
		rs.logger.Warn("extention period days is big")
		return errs.ErrExtentionPeriodDaysIsBig
	}
	if extentionPeriodDays < 3 {
		rs.logger.Warn("extention period days is small")
		return errs.ErrExtentionPeriodDaysIsSmall
	}

	rs.logger.Info("extention period days is good")

	return nil
}
