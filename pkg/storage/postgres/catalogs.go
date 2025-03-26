package postgres

import (
	"context"
	"errors"
)

// Table: catalogs

const (
	insertFuncCatalogs = "f_catalogs_insert"
	updateFuncCatalogs = "f_catalogs_update"
	deleteFuncCatalogs = "f_catalogs_delete"
	viewFuncCatalogs   = "f_catalogs_view"
)

type Catalogs struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
	DtAdd   int    `json:"dt_add"`
	Enable  bool   `json:"enable"`
}

func (s *Storage) InsertCatalogs(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogs(insertFuncCatalogs, jsonRequest)
	return id, err

}

func (s *Storage) UpdateCatalogs(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogs(updateFuncCatalogs, jsonRequest)
	return id, err

}

func (s *Storage) DeleteCatalogs(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogs(deleteFuncCatalogs, jsonRequest)
	return id, err

}

func (s *Storage) iudCatalogs(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

func (s *Storage) ViewCatalogs(jsonRequest map[string]interface{}) ([]Catalogs, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncCatalogs+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var catalogs []Catalogs

	for rows.Next() {
		var t Catalogs
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Comment,
			&t.DtAdd,
			&t.Enable,
		)
		if err != nil {
			return nil, err
		}

		catalogs = append(catalogs, t)
	}
	return catalogs, rows.Err()
}
