package teststore

import "firstRestAPI/internal/model"

type UsersLKRepository struct {
	store *Store
}

func (ulk *UsersLKRepository) Create(lk *model.UserLK, u *model.User) error {
	return nil
}
func (ulk *UsersLKRepository) FindByNickname(nick string) error {
	return nil
}
func (ulk *UsersLKRepository) LKDescriptionEdit(lk *model.UserLK) error {
	return nil
}
func (ulk *UsersLKRepository) LKShow(lk *model.UserLK) error {
	return nil
}
