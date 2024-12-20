package errs

import "errors"

var (
	ErrRatingAlreadyExist   = errors.New("[!] ratingService error! Rating already exists")
	ErrRatingDoesNotExists  = errors.New("[!] ratingService error! Rating does not exist")
	ErrRatingObjectIsNil    = errors.New("[!] ratingService error! Rating object is nil")
	ErrRatingRangeIsInvalid = errors.New("[!] ratingService error! Rating range is invalid")
)
