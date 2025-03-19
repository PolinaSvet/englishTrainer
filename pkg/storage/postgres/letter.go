package postgres

import (
	"context"
	"errors"
)

// Table: letter

const (
	insertFuncLetter = "f_letter_insert"
	updateFuncLetter = "f_letter_update"
	deleteFuncLetter = "f_letter_delete"
	viewFuncLetter   = "f_letter_view"
)

type Letter struct {
	Id   int
	Name string
}

func (s *Storage) InsertLetter(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudLetter(insertFuncLetter, jsonRequest)
	return id, err

}

func (s *Storage) UpdateLetter(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudLetter(updateFuncLetter, jsonRequest)
	return id, err

}

func (s *Storage) DeleteLetter(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudLetter(deleteFuncLetter, jsonRequest)
	return id, err

}

func (s *Storage) iudLetter(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

func (s *Storage) ViewLetter(jsonRequest map[string]interface{}) ([]Letter, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncLetter+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var letter []Letter

	for rows.Next() {
		var t Letter
		err = rows.Scan(
			&t.Id,
			&t.Name,
		)
		if err != nil {
			return nil, err
		}

		letter = append(letter, t)
	}
	return letter, rows.Err()
}
