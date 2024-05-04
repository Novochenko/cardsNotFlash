package sqlstore

import (
	"database/sql"
	"firstRestAPI/internal/store"

	_ "github.com/go-sql-driver/mysql"
)

type Store struct {
	db                *sql.DB
	userRepository    *UserRepository
	cardRepository    *CardRepository
	usersLKRepository *UsersLKRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
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
