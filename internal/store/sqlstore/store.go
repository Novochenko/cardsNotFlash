package sqlstore

import (
	"database/sql"
	"firstRestAPI/internal/store"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db                    *sql.DB
	images                string
	userRepository        *UserRepository
	cardRepository        *CardRepository
	usersLKRepository     *UsersLKRepository
	groupRepository       *GroupRepository
	cardsImagesRepository *CardsImagesRepository
}

func New(db *sql.DB, imagePath string) *Store {
	return &Store{
		db:     db,
		images: imagePath,
	}
}

func (s *Store) User() store.UserRepository {
	if s.userRepository != nil {
		return s.userRepository
	}
	s.userRepository = &UserRepository{
		store: s,
	}
	return s.userRepository
}

func (s *Store) Card() store.CardsRepository {
	if s.cardRepository != nil {
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
	return s.images
}
