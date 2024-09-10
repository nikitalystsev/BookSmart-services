package dto

type ReaderSignInDTO struct {
	PhoneNumber string `json:"phone_number"`
	Password    string `json:"password"`
}

type ReaderSignUpDTO struct {
	Fio         string `json:"fio"`
	PhoneNumber string `json:"phone_number"`
	Age         uint   `json:"age"`
	Password    string `json:"password"`
}

type ReaderTokensDTO struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiredAt    int64  `json:"expired_at"`
}
