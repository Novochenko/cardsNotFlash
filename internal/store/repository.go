package store

import "firstRestAPI/internal/model"

type UserRepository interface {
	Create(*model.User) error
	Find(int64) (*model.User, error)
	FindByEmail(string) (*model.User, error)
	Delete(id int64) error
	ShowALLGroups(u *model.User) ([]*model.Group, error)
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
	Create(*model.UserLK, *model.User) error
	FindByNickname(nick string) error
}

type GroupRepository interface {
	Create(*model.Group) error
	Delete(*model.Group) error
	Show(*model.Group) ([]*model.Card, error)
	Edit(*model.Group) error
}
