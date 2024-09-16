package dto

import "github.com/google/uuid"

type RatingInputDTO struct {
	BookID uuid.UUID `json:"book_id"`
	Review string    `json:"review"`
	Rating int       `json:"rating"`
}

type RatingOutputDTO struct {
	Reader string `json:"reader"`
	Review string `json:"review"`
	Rating int    `json:"rating"`
}

type AvgRatingDTO struct {
	AvgRating float32 `json:"avg_rating"`
}
