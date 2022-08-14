package models

type Review struct {
	Id         uint
	Reviewer   string
	MovieTitle string
	Text       string
	Rating     uint8
}
