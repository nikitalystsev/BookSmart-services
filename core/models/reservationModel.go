package models

import (
	"github.com/google/uuid"
	"time"
)

type ReservationModel struct {
	ID         uuid.UUID
	ReaderID   uuid.UUID
	BookID     uuid.UUID
	IssueDate  time.Time
	ReturnDate time.Time
	State      string
}
