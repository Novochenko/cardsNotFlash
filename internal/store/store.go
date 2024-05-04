package store

type Store interface {
	User() UserRepository
	Card() CardsRepository
	UserLK() UsersLKRepository
}
