package model

import "testing"

func TestUser(t *testing.T) *User {
	t.Helper()
	return &User{
		Email:    "user@example.org",
		Password: "password",
	}
}

func TestCard(t *testing.T) *Card {
	t.Helper()
	return &Card{
		FrontSide: "front_side",
		BackSide:  "back_side",
	}
}
