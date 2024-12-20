package dto

import (
	"github.com/google/uuid"
	"time"
)

type ReservationInputDTO struct {
	BookID uuid.UUID `json:"book_id"`
}

type ReservationOutputDTO struct {
	ID                 uuid.UUID `json:"id"`
	BookTitleAndAuthor string    `json:"book_title_and_author"`
	IssueDate          time.Time `json:"issue_date"`
	ReturnDate         time.Time `json:"return_date"`
	State              string    `json:"state"`
}

type ReservationExtentionPeriodDaysInputDTO struct {
	ExtentionPeriodDays int `json:"extention_period_days"`
}
