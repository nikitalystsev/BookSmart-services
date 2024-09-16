package models

import "github.com/google/uuid"

type RatingModel struct {
	ID       uuid.UUID
	ReaderID uuid.UUID
	BookID   uuid.UUID
	Review   string
	Rating   int
}
