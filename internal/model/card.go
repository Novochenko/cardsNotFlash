package model

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

type rawTime []byte

func (t rawTime) Time() (time.Time, error) {
	return time.Parse("15:04:05", string(t))
}

type Card struct {
	ID             int64     `json:"id"`
	UserID         int64     `json:"user_id"`
	FrontSide      string    `json:"front_side"`
	BackSide       string    `json:"back_side"`
	CardTime       time.Time `json:"card_time"`
	TimeFlag       rawTime   //`json:"time_flag"`
	TimeFlagString string    `json:"time_flag_string"`
	GroupID        int64     `json:"group_id"`
}

func (c *Card) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.FrontSide, validation.Required),
		validation.Field(&c.BackSide, validation.Required),
		validation.Field(&c.GroupID, validation.Required),
	)
}

func (c *Card) ValidateEdit() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.FrontSide, validation.Required),
		validation.Field(&c.BackSide, validation.Required),
		validation.Field(&c.ID, validation.Required),
	)
}

func (c *Card) ValidateDelete() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.ID, validation.Required),
	)
}
