package store

import "firstRestAPI/internal/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int64) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Delete(id int64) error
}

type CardsRepository interface {
	Create(*model.Card) error
	Delete(*model.Card) error
	Show(*model.Card) ([]*model.Card, error)
	Edit(*model.Card) error
	ShowUsingTime(c *model.Card) ([]*model.Card, error)
	CardFlagUp(card *model.Card) error
}

type UsersLKRepository interface {
	Create(*model.UserLK) error
	FindByNickname(nick string) error
}
