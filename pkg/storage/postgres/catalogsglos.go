package postgres

import (
	"context"
	"errors"
)

// Table: catalogsglos

const (
	insertFuncCatalogsglos = "f_catalogsglos_insert"
	updateFuncCatalogsglos = "f_catalogsglos_update"
	deleteFuncCatalogsglos = "f_catalogsglos_delete"
	viewFuncCatalogsglos   = "f_catalogsglos_view"
)

type Catalogsglos struct {
	Id         int    `json:"id"`
	GlossaryId string `json:"glossary_id"`
	CatalogsId string `json:"catalogs_id"`
	DtAdd      int    `json:"dt_add"`
}

func (s *Storage) InsertCatalogsglos(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogsglos(insertFuncCatalogsglos, jsonRequest)
	return id, err

}

func (s *Storage) UpdateCatalogsglos(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogsglos(updateFuncCatalogsglos, jsonRequest)
	return id, err

}

func (s *Storage) DeleteCatalogsglos(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudCatalogsglos(deleteFuncCatalogsglos, jsonRequest)
	return id, err

}

func (s *Storage) iudCatalogsglos(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

func (s *Storage) ViewCatalogsglos(jsonRequest map[string]interface{}) ([]Catalogsglos, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncCatalogsglos+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var catalogsglos []Catalogsglos

	for rows.Next() {
		var t Catalogsglos
		err = rows.Scan(
			&t.Id,
			&t.GlossaryId,
			&t.CatalogsId,
			&t.DtAdd,
		)
		if err != nil {
			return nil, err
		}

		catalogsglos = append(catalogsglos, t)
	}
	return catalogsglos, rows.Err()
}
