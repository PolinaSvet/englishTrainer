package postgres

import (
	"context"
	"errors"
)

// Table: mark

const (
	insertFuncMark = "f_mark_insert"
	updateFuncMark = "f_mark_update"
	deleteFuncMark = "f_mark_delete"
	viewFuncMark   = "f_mark_view"
)

type Mark struct {
	Id   int
	Name string
}

func (s *Storage) InsertMark(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudMark(insertFuncMark, jsonRequest)
	return id, err

}

func (s *Storage) UpdateMark(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudMark(updateFuncMark, jsonRequest)
	return id, err

}

func (s *Storage) DeleteMark(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudMark(deleteFuncMark, jsonRequest)
	return id, err

}

func (s *Storage) iudMark(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

	var jsonResponse SqlResponse
	err := s.db.QueryRow(context.Background(), "SELECT * FROM "+nameFunction+"($1);", jsonRequest).Scan(&jsonResponse)
	if err != nil {
		return 0, err
	}

	if jsonResponse.Err != "" {
		return 0, errors.New(jsonResponse.Err)
	}
	return jsonResponse.ID, nil
}

func (s *Storage) ViewMark(jsonRequest map[string]interface{}) ([]Mark, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncMark+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var mark []Mark

	for rows.Next() {
		var t Mark
		err = rows.Scan(
			&t.Id,
			&t.Name,
		)
		if err != nil {
			return nil, err
		}

		mark = append(mark, t)
	}
	return mark, rows.Err()
}
