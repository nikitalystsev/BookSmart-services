package impl

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"github.com/nikitalystsev/BookSmart-services/core/dto"
	"github.com/nikitalystsev/BookSmart-services/core/models"
	"github.com/nikitalystsev/BookSmart-services/errs"
	"github.com/nikitalystsev/BookSmart-services/intf"
	"github.com/nikitalystsev/BookSmart-services/intfRepo"
	"github.com/sirupsen/logrus"
)

const (
	BookRarityCommon = "Common"
	BookRarityRare   = "Rare"
	BookRarityUnique = "Unique"

	BooksPageLimit = 12
)

type BookService struct {
	bookRepo intfRepo.IBookRepo
	logger   *logrus.Entry
}

func NewBookService(bookRepo intfRepo.IBookRepo, logger *logrus.Entry) intf.IBookService {
	return &BookService{bookRepo: bookRepo, logger: logger}
}

func (bs *BookService) Create(ctx context.Context, book *models.BookModel) error {
	if book == nil {
		bs.logger.Warn("book object is nil")
		return errs.ErrBookObjectIsNil
	}

	bs.logger.Info("attempting to create book")

	if err := bs.checkBookCanBeAdded(ctx, book); err != nil {
		bs.logger.Errorf("book validation failed: %v", err)
		return err
	}

	bs.logger.Infof("creating book in repository: %+v", book)

	if err := bs.bookRepo.Create(ctx, book); err != nil {
		bs.logger.Errorf("error creating book: %v", err)
		return err
	}

	bs.logger.Info("successfully created book")

	return nil
}

func (bs *BookService) Delete(ctx context.Context, ID uuid.UUID) error {
	if ID == uuid.Nil {
		bs.logger.Warn("book object is nil")
		return errs.ErrBookObjectIsNil
	}

	bs.logger.Infof("attempting to delete book with ID: %s", ID)

	existingBook, err := bs.bookRepo.GetByID(ctx, ID)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		bs.logger.Errorf("error checking book existence: %v", err)
		return err
	}

	if existingBook == nil {
		bs.logger.Warn("book with this ID does not exist")
		return errs.ErrBookDoesNotExists
	}

	if err = bs.bookRepo.Delete(ctx, ID); err != nil {
		bs.logger.Errorf("error deleting book with ID %s: %v", ID, err)
		return err
	}

	bs.logger.Infof("successfully deleted book with ID: %s", ID)

	return nil
}

func (bs *BookService) GetByID(ctx context.Context, ID uuid.UUID) (*models.BookModel, error) {
	bs.logger.Infof("attempting to get book with ID: %s", ID)

	book, err := bs.bookRepo.GetByID(ctx, ID)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		bs.logger.Errorf("error checking book existence: %v", err)
		return nil, err
	}

	if book == nil {
		bs.logger.Warn("book with this ID does not exist")
		return nil, errs.ErrBookDoesNotExists
	}

	bs.logger.Infof("successfully getting book by ID: %s", ID)

	return book, nil
}

func (bs *BookService) GetByParams(ctx context.Context, params *dto.BookParamsDTO) ([]*models.BookModel, error) {
	bs.logger.Infof("attempting to search for books with params")

	books, err := bs.bookRepo.GetByParams(ctx, params)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		bs.logger.Errorf("error searching books with params: %v", err)
		return nil, err
	}

	if books == nil {
		bs.logger.Warn("books with params does not exist")
		return nil, errs.ErrBookDoesNotExists
	}

	bs.logger.Infof("successfully found %d books with params", len(books))

	return books, nil
}

func (bs *BookService) checkBookCanBeAdded(ctx context.Context, book *models.BookModel) error {
	existingBook, err := bs.bookRepo.GetByID(ctx, book.ID)
	if err != nil && !errors.Is(err, errs.ErrBookDoesNotExists) {
		bs.logger.Errorf("error checking book existence: %v", err)
		return err
	}

	if existingBook != nil {
		bs.logger.Warn("book with this ID already exists")
		return errs.ErrBookAlreadyExist
	}

	if book.Title == "" {
		bs.logger.Warn("empty book title")
		return errs.ErrEmptyBookTitle
	}

	if book.Author == "" {
		bs.logger.Warn("empty book author")
		return errs.ErrEmptyBookAuthor
	}

	if book.Rarity == "" {
		bs.logger.Warn("empty book rarity")
		return errs.ErrEmptyBookRarity
	}

	if book.CopiesNumber <= 0 {
		bs.logger.Warn("invalid book copies number")
		return errs.ErrInvalidBookCopiesNum
	}

	bs.logger.Info("book validation successful")

	return nil
}
