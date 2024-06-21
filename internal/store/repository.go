package store

import (
	"firstRestAPI/internal/model"
	"os"

	"github.com/google/uuid"
)

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
	LKDescriptionEdit(lk *model.UserLK) error
	LKShow(lk *model.UserLK) error
}

type GroupRepository interface {
	Create(*model.Group) error
	Delete(*model.Group) error
	Show(*model.Group) ([]*model.Card, error)
	Edit(*model.Group) error
	ShowUsingTime(g *model.Group) ([]*model.Card, error)
}

type CardsImagesRepository interface {
	Add(image *model.CardImages, isFrontSide bool) (*os.File, error)
	CardIDUpdate(cardID int64, uuID uuid.UUID) error
	Delete(images []*model.CardImages) error
	DeleteOnError(imageIDs uuid.UUIDs) error
	Edit(images []*model.CardImages, buf [][]byte) error
	Show(cards []*model.Card) ([]*model.CardImages, error)
}
