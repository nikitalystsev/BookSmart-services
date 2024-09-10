package dto

type BookParamsDTO struct {
	Title          string
	Author         string
	Publisher      string
	CopiesNumber   uint
	Rarity         string
	Genre          string
	PublishingYear uint
	Language       string
	AgeLimit       uint
	Limit          uint
	Offset         int
}

type BookDTO struct {
	Title          string `json:"title"`
	Author         string `json:"author"`
	Publisher      string `json:"publisher"`
	CopiesNumber   uint   `json:"copies_number"`
	Rarity         string `json:"rarity"`
	Genre          string `json:"genre"`
	PublishingYear uint   `json:"publishing_year"`
	Language       string `json:"language"`
	AgeLimit       uint   `json:"age_limit"`
}
