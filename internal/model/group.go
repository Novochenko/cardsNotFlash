package model

import (
	validation "github.com/go-ozzo/ozzo-validation"
)

type Group struct {
	GroupName string `json:"group_name"`
	GroupID   int64  `json:"group_id"`
	UserID    int64  `json:"user_id"`
}

func (g *Group) Validate() error {
	return validation.ValidateStruct(
		g,
		validation.Field(&g.GroupName, validation.Required),
	)
}
func (g *Group) ValidateEdit() error {
	return validation.ValidateStruct(
		g,
		validation.Field(&g.GroupName, validation.Required),
		validation.Field(&g.GroupID, validation.Required),
	)
}

func (g *Group) ValidateDelete() error {
	return validation.ValidateStruct(
		g,
		validation.Field(&g.GroupID, validation.Required),
	)
}
