package errs

import "errors"

var (
	ErrBookAlreadyExist      = errors.New("[!] bookService error! Book already exists")
	ErrBookAlreadyIsFavorite = errors.New("[!] bookService error! Book already in favorites")
	ErrEmptyBookTitle        = errors.New("[!] bookService error! Empty book title")
	ErrEmptyBookAuthor       = errors.New("[!] bookService error! Empty book author")
	ErrEmptyBookRarity       = errors.New("[!] bookService error! Empty book rarity")
	ErrInvalidBookCopiesNum  = errors.New("[!] bookService error! Invalid book copies number")
	ErrBookDoesNotExists     = errors.New("[!] bookService error! Book does not exist")
	ErrBookNoCopiesNum       = errors.New("[!] bookService error! Book no copies number")
	ErrBookObjectIsNil       = errors.New("[!] bookService error! Book object is nil")
)
