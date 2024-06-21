package sqlstore

import (
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
)

type UsersLKRepository struct {
	store *Store
}

func (ulk *UsersLKRepository) Create(lk *model.UserLK, u *model.User) error {
	if err := lk.Validate(); err != nil {
		return err
	}
	stmt, err := ulk.store.db.Prepare("INSERT lk (user_id, nickname, email, encrypted_password) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()
	res, err := stmt.Exec(lk.UserID, lk.Nickname, u.Email, u.EncryptedPassword)
	if err != nil {
		return err
	}
	lk.ID, err = res.LastInsertId()
	if err != nil {
		return err
	}
	return nil
}

// если запись есть, возвращает ошибку
func (ulk *UsersLKRepository) FindByNickname(nick string) error {
	var exists bool
	row := ulk.store.db.QueryRow("SELECT EXISTS(SELECT user_id FROM lk WHERE nickname = ?)", nick)
	if err := row.Scan(&exists); err != nil {
		return err
	}
	if exists {
		return store.ErrNicknameNotAvialable
	}
	return nil
}
func (ulk *UsersLKRepository) LKShow(lk *model.UserLK) error {
	if err := lk.ValidateShow(); err != nil {
		return err
	}
	row := ulk.store.db.QueryRow("SELECT email, nickname, cards_count, user_description FROM lk WHERE user_id = ?", lk.UserID)
	if err := row.Scan(&lk.Email, &lk.Nickname, &lk.CardsCount, &lk.UserDescription); err != nil {
		return err
	}
	return nil
}

func (ulk *UsersLKRepository) LKDescriptionEdit(lk *model.UserLK) error {
	if lk == nil {
		return store.ErrEmptyLK
	}
	stmt, err := ulk.store.db.Prepare("UPDATE lk SET user_description = ? WHERE user_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	if _, err = stmt.Exec(lk.UserDescription, lk.UserID); err != nil {
		return err
	}
	return nil
}
