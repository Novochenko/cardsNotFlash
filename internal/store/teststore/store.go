package teststore

import (
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
)

type Store struct {
	userRepository        *UserRepository
	cardRepository        *CardRepository
	usersLKRepository     *UsersLKRepository
	groupRepository       *GroupRepository
	cardsImagesRepository *CardsImagesRepository
}

// New ...
func New() *Store {
	return &Store{}
}

// User ...
func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}

	s.userRepository = &UserRepository{
		store: s,
		users: make(map[int64]*model.User),
	}

	return s.userRepository
}

func (s *Store) Card() store.CardsRepository {
	if s.userRepository != nil {
		return s.cardRepository
	}
	s.cardRepository = &CardRepository{
		store: s,
	}
	return s.cardRepository
}

func (s *Store) UserLK() store.UsersLKRepository {
	if s.usersLKRepository != nil {
		return s.usersLKRepository
	}
	s.usersLKRepository = &UsersLKRepository{
		store: s,
	}
	return s.usersLKRepository
}
func (s *Store) Group() store.GroupRepository {
	if s.groupRepository != nil {
		return s.groupRepository
	}
	s.groupRepository = &GroupRepository{
		store: s,
	}
	return s.groupRepository
}

func (s *Store) CardImages() store.CardsImagesRepository {
	if s.cardsImagesRepository != nil {
		return s.cardsImagesRepository
	}
	s.cardsImagesRepository = &CardsImagesRepository{
		store: s,
	}
	return s.cardsImagesRepository
}

func (s *Store) Images() string {
	return " "
}
