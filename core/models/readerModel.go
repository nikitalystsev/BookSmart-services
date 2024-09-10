package models

import "github.com/google/uuid"

type ReaderModel struct {
	ID          uuid.UUID
	Fio         string
	PhoneNumber string
	Age         uint
	Password    string
	Role        string
}

type Tokens struct {
	AccessToken  string
	RefreshToken string
}
