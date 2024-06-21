package sqlstore

import (
	"bytes"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
	"io"
	"os"

	"github.com/google/uuid"
)

type CardsImagesRepository struct {
	store *Store
}

func (cir *CardsImagesRepository) Show(cards []*model.Card) ([]*model.CardImages, error) {
	cardImages := make([]*model.CardImages, len(cards))
	for _, v := range cards {
		rows, err := cir.store.db.Query("SELECT image_id FROM card_images WHERE card_id = ?", v.ID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			cardImg := &model.CardImages{
				CardID: v.ID,
			}
			if err = rows.Scan(&cardImg.ImageID); err != nil {
				return nil, err
			}
		}
	}
	return cardImages, nil
}

func (cir *CardsImagesRepository) addHelper(image *model.CardImages, isFrontSide bool) (*os.File, error) {
	stmt, err := cir.store.db.Prepare(`INSERT card_images(image_id, user_id, front_side, back_side) 
								   VALUES (UUID_TO_BIN(?), ?, ?, ?)`)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var IsBackSide bool
	if !isFrontSide {
		isFrontSide = true
	}
	_, err = stmt.Exec(image.ImageID.String(), image.UserID, isFrontSide, IsBackSide)
	if err != nil {
		return nil, err
	}
	file, err := os.Create(cir.store.Images() + "/" + "cardimages" + "/" + image.ImageID.String() + ".png")
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (cir *CardsImagesRepository) Add(image *model.CardImages, isFrontSide bool) (*os.File, error) {
	file, err := cir.addHelper(image, isFrontSide)
	if err != nil {
		return nil, err
	}
	return file, nil
}
func (cir *CardsImagesRepository) CardIDUpdate(cardID int64, uuID uuid.UUID) error {
	stmt, err := cir.store.db.Prepare(`UPDATE card_images
									   SET card_id = ?
									   WHERE image_id = UUID_TO_BIN(?);`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(cardID, uuID.String())
	if err != nil {
		return err
	}
	return nil
}
func (cir *CardsImagesRepository) deleteHelper(image *model.CardImages) error {
	stmt, err := cir.store.db.Prepare("DELETE FROM card_images WHERE image_id = UUID_TO_BIN(?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(image.ImageID); err != nil {
		return err
	}

	if err = os.Remove(cir.store.Images() + "/" + "cardimages" + "/" + image.ImageID.String() + ".png"); err != nil {
		return err
	}
	return nil
}

func (cir *CardsImagesRepository) Delete(images []*model.CardImages) error {
	for k := range images {
		if err := cir.deleteHelper(images[k]); err != nil {
			return err
		}
	}
	return nil
}
func (cir *CardsImagesRepository) DeleteOnCascade(card *model.Card) error {
	rows, err := cir.store.db.Query("SELECT BIN_TO_UUID(image_id) WHERE card_id = ?", card.ID)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var fileName string
		if err := rows.Scan(&fileName); err != nil {
			return err
		}
		stmt, err := cir.store.db.Prepare(`DELETE FROM card_images
									   WHERE image_id = UUID_TO_BIN(?) `)
		if err != nil {
			return err
		}
		if _, err = stmt.Exec(fileName); err != nil {
			return err
		}
		if err := os.Remove(cir.store.Images() + "/" + "cardimages" + "/" + fileName + ".png"); err != nil {
			return err
		}
	}
	return nil
}

func (cir *CardsImagesRepository) DeleteOnError(imageIDs uuid.UUIDs) error {
	for _, v := range imageIDs {
		stmt, err := cir.store.db.Prepare(`DELETE FROM card_images
									   WHERE image_id = UUID_TO_BIN(?) `)
		if err != nil {
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(v.String()); err != nil {
			return err
		}
		if err := os.Remove(cir.store.Images() + "/" + "cardimages" + "/" + v.String() + ".png"); err != nil {
			return err
		}
	}
	return nil
}

func (cir *CardsImagesRepository) editHelper(image *model.CardImages, buf []byte) error {

	path := cir.store.Images() + "/" + "cardimages" + "/" + image.ImageID.String() + ".png"
	if err := os.Remove(path); err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	buffer := bytes.NewBuffer(buf)
	if _, err := io.Copy(file, buffer); err != nil {
		return err
	}
	return nil
}

func (cir *CardsImagesRepository) Edit(images []*model.CardImages, buf [][]byte) error {
	if len(images) != len(buf) {
		return store.ErrNotMatch
	}
	for k := range images {
		cir.editHelper(images[k], buf[k])
	}
	return nil
}
