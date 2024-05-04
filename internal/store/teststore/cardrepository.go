package teststore

import (
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
)

type CardRepository struct {
	store *Store
	//
	userCards map[int64]map[int64]*model.Card
}

func (cr *CardRepository) ShowUsingTime(c *model.Card) ([]*model.Card, error) {
	d := []*model.Card{}
	return d, nil
}

func (cr *CardRepository) Create(c *model.Card) error {
	ID := int64(len(cr.userCards[c.UserID]))
	c.ID = ID
	if a := cr.userCards[c.UserID]; a == nil {
		cr.userCards[c.UserID] = make(map[int64]*model.Card)
	}
	cr.userCards[c.UserID][ID] = c

	return nil
}
func (cr *CardRepository) Delete(c *model.Card) error {
	if a := cr.userCards[c.UserID]; a == nil {
		return store.ErrRecordNotFound
	}
	if _, ok := cr.userCards[c.UserID][c.ID]; ok {
		delete(cr.userCards[c.UserID], c.ID)
	} else {
		return store.ErrRecordNotFound
	}
	return nil
}
func (cr *CardRepository) Show(c *model.Card) ([]*model.Card, error) {
	return nil, nil
}

func (cr *CardRepository) Edit(c *model.Card) error {
	return nil
}
func (cr *CardRepository) CardFlagUp(card *model.Card) error {
	return nil
}
