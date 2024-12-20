package dto

import "github.com/google/uuid"

type SignInInputDTO struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type SignUpInputDTO struct {
	Fio         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	Age         uint   `json:"age"`
	Password    string `json:"password"`
}

type SignInOutputDTO struct {
	ReaderID     uuid.UUID `json:"reader_id"`
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	ExpiredAt    int64     `json:"expired_at"`
}

type RefreshTokenInputDTO struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutputDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
}

type FavoriteBookInputDTO struct {
	BookID uuid.UUID `json:"book_id"`
}
