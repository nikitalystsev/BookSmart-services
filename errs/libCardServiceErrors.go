package errs

import "errors"

var (
	ErrLibCardAlreadyExist  = errors.New("[!] libCardService error! LibCard already exists")
	ErrLibCardDoesNotExists = errors.New("[!] libCardService error! LibCard does not exist")
	ErrLibCardIsValid       = errors.New("[!] libCardService error! LibCard is valid")
	ErrLibCardIsInvalid     = errors.New("[!] libCardService error! LibCard is invalid")
	ErrLibCardObjectIsNil   = errors.New("[!] libCardService error! LibCard object is nil")
)
