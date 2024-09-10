package models

import (
	"github.com/google/uuid"
	"time"
)

type LibCardModel struct {
	ID           uuid.UUID
	ReaderID     uuid.UUID
	LibCardNum   string
	Validity     int
	IssueDate    time.Time
	ActionStatus bool
}
