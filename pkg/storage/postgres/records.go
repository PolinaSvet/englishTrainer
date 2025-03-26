package postgres

import (
	"context"
	"errors"
)

// Table: records

const (
	insertFuncRecords = "f_records_insert"
	updateFuncRecords = "f_records_update"
	deleteFuncRecords = "f_records_delete"
	viewFuncRecords   = "f_records_view"
)

type Records struct {
	Id         int `json:"id"`
	UsersId    int `json:"users_id"`
	GlossaryId int `json:"glossary_id"`
	CatalogsId int `json:"catalogs_id"`
	Attempt    int `json:"attempt"`
	Guess      int `json:"guess"`
	DtAdd      int `json:"dt_add"`
}

func (s *Storage) InsertRecords(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudRecords(insertFuncRecords, jsonRequest)
	return id, err

}

func (s *Storage) UpdateRecords(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudRecords(updateFuncRecords, jsonRequest)
	return id, err

}

func (s *Storage) DeleteRecords(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudRecords(deleteFuncRecords, jsonRequest)
	return id, err

}

func (s *Storage) iudRecords(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

func (s *Storage) ViewRecords(jsonRequest map[string]interface{}) ([]Records, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncRecords+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var records []Records

	for rows.Next() {
		var t Records
		err = rows.Scan(
			&t.Id,
			&t.UsersId,
			&t.GlossaryId,
			&t.CatalogsId,
			&t.Attempt,
			&t.Guess,
			&t.DtAdd,
		)
		if err != nil {
			return nil, err
		}

		records = append(records, t)
	}
	return records, rows.Err()
}
