package sqlstore

import (
	"database/sql"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
	"time"
)

const (
	Start uint = iota
	twenty_m
	one_h
	two_h
	four_h
	nine_h
	one_d
	End
)

var (
	cardFlags = map[uint]string{
		Start:    "00:00:01",
		twenty_m: "00:20:00",
		one_h:    "01:00:00",
		two_h:    "02:00:00",
		four_h:   "04:00:00",
		nine_h:   "09:00:00",
		one_d:    "23:59:59",
	}
)

type CardRepository struct {
	store *Store
}

func (cr *CardRepository) Create(c *model.Card) error {
	if err := c.Validate(); err != nil {
		return err
	}
	stmt, err := cr.store.db.Prepare("INSERT cards (user_id, front_side, back_side, card_time, time_flag) VALUES(?, ?, ?, CURRENT_TIMESTAMP(), ?);")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(c.UserID, c.FrontSide, c.BackSide, Start)
	if err != nil {
		return err
	}
	if c.ID, err = res.LastInsertId(); err != nil {
		return err
	}
	return nil
}
func (cr *CardRepository) Delete(c *model.Card) error {
	if err := c.ValidateDelete(); err != nil {
		return err
	}
	stmt, err := cr.store.db.Prepare("DELETE FROM cards WHERE card_id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.ID, c.UserID)
	if err != nil {
		return err
	}

	return nil
}

func (cr *CardRepository) Show(c *model.Card) ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := cr.store.db.Query("SELECT card_id, user_id, front_side, back_side, card_time, time_flag FROM cards WHERE user_id = ?;", c.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return cards, store.ErrRecordNotFound
		}
		return cards, err
	}
	for rows.Next() {
		scans := &model.Card{}
		if err := rows.Scan(&scans.ID, &scans.UserID, &scans.FrontSide, &scans.BackSide, &scans.CardTime, &scans.TimeFlag); err != nil {
			return nil, err
		}
		cards = append(cards, scans)
	}
	if len(cards) == 0 {
		return cards, store.ErrRecordNotFound
	}
	return cards, nil
}

func (cr *CardRepository) Edit(c *model.Card) error {
	if err := c.ValidateEdit(); err != nil {
		return err
	}
	stmt, err := cr.store.db.Prepare(`UPDATE cards
						SET front_side = ?, back_side = ?, card_time = current_timestamp(), time_flag = ? WHERE card_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(c.FrontSide, c.BackSide, Start, c.ID, c.UserID)
	if err != nil {
		return err
	}
	// if _, err := res.RowsAffected(); err != nil{
	// 	return err
	// }

	return nil
}

/*
сделай так, чтобы через один запрос в бд приходили карты,
то есть чтобы автоматом флаг превращался в минуты и попадал в селект запрос
*/
func (cr *CardRepository) ShowUsingTime(c *model.Card) ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := cr.store.db.Query(`SELECT card_id, front_side, back_side, time_flag, card_time, time_flag
									FROM cards
								 	WHERE user_id = ? AND TIMEDIFF(NOW(), card_time) > time_flag;
								 `, c.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return cards, store.ErrRecordNotFound
		}
		return cards, err
	}
	for rows.Next() {
		scans := &model.Card{}
		var ct time.Time
		if err := rows.Scan(&scans.ID, &scans.FrontSide, &scans.BackSide, &scans.TimeFlag, &ct, &scans.TimeFlagString); err != nil {
			return nil, err
		}
		scans.CardTime = ct
		cards = append(cards, scans)
	}
	if len(cards) == 0 {
		return cards, store.ErrRecordNotFound
	}
	// if err = cr.cardFlagUp(cards); err != nil {
	// 	return nil, err
	// }
	return cards, nil
}

// проблема с работой глобальной переменной!!!
func (cr *CardRepository) CardFlagUp(card *model.Card) error {
	row := cr.store.db.QueryRow(`SELECT time_flag
								 FROM cards
	 							 WHERE user_id = ? AND card_id = ?;`, card.UserID, card.ID)
	if err := row.Scan(&card.TimeFlag); err != nil {
		return err
	}
	tmpTime, err := card.TimeFlag.Time()
	if err != nil {
		return err
	}
	var currentT string

	for i := 0; i < len(cardFlags); i++ {
		v := cardFlags[uint(i)]
		v1, err := time.Parse("15:04:05", v)
		if err != nil {
			return err
		}
		if tmpTime.Before(v1) {
			currentT = v
			break
		}
	}
	stmt, err := cr.store.db.Prepare(`UPDATE cards
										SET time_flag = ?, card_time = current_timestamp()
										WHERE card_id = ?;`)
	if err != nil {
		return err
	}
	if _, err = stmt.Exec(currentT, card.ID); err != nil {
		return err
	}
	row = cr.store.db.QueryRow("SELECT card_time FROM cards WHERE card_id = ?", card.ID)
	if err = row.Scan(&card.CardTime); err != nil {
		return err
	}
	card.TimeFlagString = currentT

	return nil
}
