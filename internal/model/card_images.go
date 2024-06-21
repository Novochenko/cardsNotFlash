package model

import "github.com/google/uuid"

// import (

// 	validation "github.com/go-ozzo/ozzo-validation"
// )

type CardImages struct {
	ImageID     uuid.UUID `json:"image_id"`
	CardID      int64     `json:"card_id"`
	UserID      int64     `json:"user_id"`
	IsFrontSide bool      `json:"is_front_side"`
	IsBackSide  bool      `json:"is_back_side"`
}

// func (c *CardImages) Validate() error {
// 	return validation.ValidateStruct(
// 		c,
// 		validation.Field(&c.FrontSide, validation.Required),
// 		validation.Field(&c.BackSide, validation.Required),
// 		validation.Field(&c.GroupID, validation.Required),
// 	)
// }
