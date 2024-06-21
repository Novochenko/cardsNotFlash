package teststore

import (
	"firstRestAPI/internal/model"
	"os"

	"github.com/google/uuid"
)

type CardsImagesRepository struct {
	store *Store
}

func (cir *CardsImagesRepository) Add(image *model.CardImages, isFrontSide bool) (*os.File, error) {
	return nil, nil
}

func (cir *CardsImagesRepository) Delete(images []*model.CardImages) error {
	return nil
}

func (cir *CardsImagesRepository) Edit(images []*model.CardImages, buf [][]byte) error {
	return nil
}
func (cir *CardsImagesRepository) Show(cards []*model.Card) ([]*model.CardImages, error) {
	return nil, nil
}
func (cir *CardsImagesRepository) CardIDUpdate(cardID int64, uuID uuid.UUID) error {
	return nil
}
func (cir *CardsImagesRepository) DeleteOnError(imageIDs uuid.UUIDs) error {
	return nil
}
