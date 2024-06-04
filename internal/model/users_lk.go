package model

import (
	"os"

	validation "github.com/go-ozzo/ozzo-validation"
)

type UserLK struct {
	ID              int64    `json:"id" structs:"id"`
	UserID          int64    `json:"user_id" structs:"user_id"`
	Nickname        string   `json:"nickname" structs:"nickname"`
	CardsCount      int64    `json:"cards_count" structs:"cards_count"`
	UserDescription string   `json:"user_desription" structs:"user_desription"`
	Email           string   `json:"email" structs:"email"`
	Image           *os.File `json:"image" structs:"image"`
}

func (ulk *UserLK) Validate() error {
	return validation.ValidateStruct(
		ulk,
		validation.Field(&ulk.Nickname, validation.Required))
}
func (ulk *UserLK) ValidateShow() error {
	return validation.ValidateStruct(
		ulk,
		validation.Field(&ulk.UserID, validation.Required),
	)
}
