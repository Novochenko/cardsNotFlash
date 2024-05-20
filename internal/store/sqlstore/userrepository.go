package sqlstore

import (
	"database/sql"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
)

type UserRepository struct {
	store *Store
}

func (ur *UserRepository) ShowALLGroups(u *model.User) ([]*model.Group, error) {
	groups := []*model.Group{}
	rows, err := ur.store.db.Query("SELECT group_id, group_name, user_id FROM card_groups WHERE user_id = ?", u.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	for rows.Next() {
		g := &model.Group{}
		rows.Scan(&g.GroupID, &g.GroupName, &g.UserID)
		groups = append(groups, g)
	}
	return groups, nil
}

func (ur *UserRepository) Create(u *model.User) error {
	if err := u.EncryptPassword(); err != nil {
		return err
	}

	stmt, err := ur.store.db.Prepare("INSERT users (email, encrypted_password) VALUES (?, ?);")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(u.Email, u.EncryptedPassword)
	if err != nil {
		return err
	}
	stmt.Close() // чек такой
	if u.ID, err = res.LastInsertId(); err != nil {
		return err
	}

	return nil
}
func (ur *UserRepository) Delete(id int64) error {
	stmt, err := ur.store.db.Prepare("DELETE from users WHERE user_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
func (ur *UserRepository) FindByEmail(email string) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM users WHERE email = ?", email,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}

func (ur *UserRepository) Find(id int64) (*model.User, error) {
	u := &model.User{}
	if err := ur.store.db.QueryRow(
		"SELECT user_id, email, encrypted_password FROM users WHERE user_id = ?", id,
	).Scan(
		&u.ID,
		&u.Email,
		&u.EncryptedPassword,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}

	return u, nil
}
