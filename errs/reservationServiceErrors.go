package errs

import "errors"

var (
	ErrReservationsLimitExceeded    = errors.New("[!] reservationService error! Reservations limit exceeded")
	ErrUniqueBookNotReserved        = errors.New("[!] reservationService error! Unique book not reserved")
	ErrRareAndUniqueBookNotExtended = errors.New("[!] reservationService error! Rare and unique book not extended")
	ErrReservationAgeLimit          = errors.New("[!] reservationService error! The reader's age is less than the limit")
	ErrReservationIsAlreadyClosed   = errors.New("[!] reservationService error! Reservation is already closed")
	ErrReservationIsAlreadyExpired  = errors.New("[!] reservationService error! Reservation is already expired")
	ErrReservationIsAlreadyExtended = errors.New("[!] reservationService error! Reservation is already extended")
	ErrReservationObjectIsNil       = errors.New("[!] reservationService error! Reservation object is nil")
	ErrReservationDoesNotExists     = errors.New("[!] reservationService error! Reservation does not exists")
	ErrReservationAlreadyExists     = errors.New("[!] reservationService error! Reservation already exists")
)
