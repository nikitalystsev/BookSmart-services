package dto

import "github.com/google/uuid"

type RatingInputDTO struct {
	ReaderID uuid.UUID `json:"reader_id"`
	Review   string    `json:"review"`
	Rating   int       `json:"rating"`
}

type RatingOutputDTO struct {
	ReaderFio string `json:"reader_fio"`
	Review    string `json:"review"`
	Rating    int    `json:"rating"`
}

type AvgRatingOutputDTO struct {
	AvgRating float32 `json:"avg_rating"`
}
