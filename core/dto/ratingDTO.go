package dto

import "github.com/google/uuid"

type RatingDTO struct {
	BookID uuid.UUID `json:"book_id"`
	Review string    `json:"review"`
	Rating int       `json:"rating"`
}
