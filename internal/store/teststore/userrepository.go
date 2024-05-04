package teststore

import (
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
)

type UserRepository struct {
	store *Store
	users map[int64]*model.User
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return err
	}
	if err := u.EncryptPassword(); err != nil {
		return err
	}

	u.ID = int64(len(ur.users))
	ur.users[int64(u.ID)] = u
	return nil
}
func (ur *UserRepository) Delete(id int64) error {
	return nil
}

func (ur *UserRepository) Find(id int64) (*model.User, error) {
	if u, ok := ur.users[id]; ok {
		return u, nil
	} else {
		return nil, store.ErrRecordNotFound
	}
}

func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	for _, u := range ur.users {
		if u.Email == email {
			return u, nil
		}
	}
	return nil, store.ErrRecordNotFound
}
