package store

type Store interface {
	User() UserRepository
	Card() CardsRepository
	UserLK() UsersLKRepository
	Group() GroupRepository
	CardImages() CardsImagesRepository
	Images() string
}
