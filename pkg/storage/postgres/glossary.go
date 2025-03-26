package postgres

import (
	"context"
	"encoding/json"
	"errors"
)

// Table: glossary

const (
	insertFuncGlossary     = "f_glossary_insert"
	updateFuncGlossary     = "f_glossary_update"
	deleteFuncGlossary     = "f_glossary_delete"
	viewFuncGlossary       = "f_glossary_view"
	viewRandomFuncGlossary = "f_glossary_view_random"
)

type ExampleColumn struct {
	Ex1 string `json:"ex1"`
	Ex2 string `json:"ex2"`
}

type Glossary struct {
	Id            int             `json:"id"`
	Letter        string          `json:"letter"`
	Word          string          `json:"word"`
	Transcription string          `json:"transcription"`
	Translation   string          `json:"translation"`
	Example       []ExampleColumn `json:"example"`
	DtAdd         int             `json:"dt_add"`
	Enable        bool            `json:"enable"`
	CatalogId     int             `json:"catalog_id"`
	UsersId       int             `json:"users_id"`
	RecordsId     int             `json:"records_id"`
	Attempt       int             `json:"attempt"`
	Guess         bool            `json:"guess"`
}

// 1. uid
func (s *Storage) InsertGlossary(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudGlossary(insertFuncGlossary, jsonRequest)
	return id, err
}

func (s *Storage) UpdateGlossary(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudGlossary(updateFuncGlossary, jsonRequest)
	return id, err
}

func (s *Storage) DeleteGlossary(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudGlossary(deleteFuncGlossary, jsonRequest)
	return id, err
}

func (s *Storage) iudGlossary(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

// 2. select
func (s *Storage) ViewGlossary(jsonRequest map[string]interface{}) ([]Glossary, error) {

	data, err := s.viewGlossary(viewFuncGlossary, jsonRequest)
	return data, err
}

func (s *Storage) ViewRandomGlossary(jsonRequest map[string]interface{}) ([]map[string]interface{}, error) {

	data, err := s.viewGlossary(viewRandomFuncGlossary, jsonRequest)

	var result []map[string]interface{}

	for _, glossary := range data {
		// Конвертируем поле Example в []map[string]interface{}
		exampleMaps := make([]map[string]interface{}, len(glossary.Example))
		for i, example := range glossary.Example {
			exampleMaps[i] = map[string]interface{}{
				"ex1": example.Ex1,
				"ex2": example.Ex2,
			}
		}

		// Создаем map для текущего Glossary
		m := map[string]interface{}{
			"id":            glossary.Id,
			"letter":        glossary.Letter,
			"word":          glossary.Word,
			"transcription": glossary.Transcription,
			"translation":   glossary.Translation,
			"example":       exampleMaps,
			"dt_add":        glossary.DtAdd,
			"enable":        glossary.Enable,
			"catalog_id":    glossary.CatalogId,
			"users_id":      glossary.UsersId,
			"records_id":    glossary.RecordsId,
			"attempt":       glossary.Attempt,
			"guess":         glossary.Guess,
		}
		result = append(result, m)
	}

	return result, err
}

func (s *Storage) viewGlossary(nameFunction string, jsonRequest map[string]interface{}) ([]Glossary, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+nameFunction+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var glossary []Glossary

	for rows.Next() {
		var t Glossary
		var exampleData []byte

		err = rows.Scan(
			&t.Id,
			&t.Letter,
			&t.Word,
			&t.Transcription,
			&t.Translation,
			&exampleData,
			&t.DtAdd,
			&t.Enable,
			&t.CatalogId,
			&t.UsersId,
			&t.RecordsId,
			&t.Attempt,
			&t.Guess,
		)
		if err != nil {
			return nil, err
		}

		if err = json.Unmarshal(exampleData, &t.Example); err != nil {
			return nil, err
		}

		glossary = append(glossary, t)
	}
	return glossary, rows.Err()
}
