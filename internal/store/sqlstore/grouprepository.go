package sqlstore

import (
	"database/sql"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
	"time"
)

type GroupRepository struct {
	store *Store
}

func (gr *GroupRepository) Create(g *model.Group) error {
	if err := g.Validate(); err != nil {
		return err
	}
	stmt, err := gr.store.db.Prepare("INSERT card_groups (group_name, user_id) VALUES(?, ?);")
	if err != nil {
		return err
	}
	res, err := stmt.Exec(g.GroupName, g.UserID)
	if err != nil {
		return err
	}
	if g.GroupID, err = res.LastInsertId(); err != nil {
		return err
	}
	return nil
}
func (gr *GroupRepository) Delete(g *model.Group) error {
	if err := g.ValidateDelete(); err != nil {
		return err
	}
	stmt, err := gr.store.db.Prepare("DELETE FROM card_groups WHERE group_id = ? AND user_id = ?")
	if err != nil {
		return err
	}
	_, err = stmt.Exec(g.GroupID, g.UserID)
	if err != nil {
		return err
	}
	return nil
}
func (gr *GroupRepository) Show(g *model.Group) ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := gr.store.db.Query("SELECT card_id, user_id, front_side, back_side, card_time, time_flag FROM cards WHERE user_id = ? AND group_id = ?;", g.UserID, g.GroupID)
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
func (gr *GroupRepository) ShowUsingTime(g *model.Group) ([]*model.Card, error) {
	cards := []*model.Card{}
	rows, err := gr.store.db.Query(`SELECT card_id, front_side, back_side, time_flag, card_time, time_flag
									FROM cards
								 	WHERE user_id = ? AND group_id = ? AND TIMEDIFF(NOW(), card_time) > time_flag;
								 `, g.UserID, g.GroupID)
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
	return cards, nil
}
func (gr *GroupRepository) Edit(g *model.Group) error {
	if err := g.ValidateEdit(); err != nil {
		return err
	}
	stmt, err := gr.store.db.Prepare(`UPDATE card_groups
						SET group_name = ? WHERE group_id = ? AND user_id = ?;`)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(g.GroupName, g.GroupID, g.UserID)
	if err != nil {
		return err
	}
	return nil
}
