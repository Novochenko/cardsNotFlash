package teststore

import "firstRestAPI/internal/model"

type GroupRepository struct {
	store *Store
}

func (gr *GroupRepository) Create(g *model.Group) error {
	return nil
}
func (gr *GroupRepository) Delete(g *model.Group) error {
	return nil
}
func (gr *GroupRepository) Show(g *model.Group) ([]*model.Card, error) {
	var cards []*model.Card
	return cards, nil
}
func (gr *GroupRepository) Edit(g *model.Group) error {
	return nil
}
func (gr *GroupRepository) ShowUsingTime(g *model.Group) ([]*model.Card, error) {
	return nil, nil
}
