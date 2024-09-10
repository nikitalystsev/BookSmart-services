package models

import "github.com/google/uuid"

type BookModel struct {
	ID             uuid.UUID
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   uint
	Rarity         string
	Genre          string
	PublishingYear uint
	Language       string
	AgeLimit       uint
}
