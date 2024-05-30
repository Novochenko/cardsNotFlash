package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type UserLK struct {
	ID              int64  `json:"id"`
	UserID          int64  `json:"user_id"`
	Nickname        string `json:"nickname"`
	CardsCount      int64  `json:"cards_count"`
	UserDescription string `json:"user_desription"`
	Email           string `json:"email"`
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
