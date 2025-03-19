package postgres

import (
	"context"
	"errors"
)

// Table: users

const (
	insertFuncUsers = "f_users_insert"
	updateFuncUsers = "f_users_update"
	deleteFuncUsers = "f_users_delete"
	viewFuncUsers   = "f_users_view"
)

type Users struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Mail     string `json:"mail"`
	Telega   string `json:"telega"`
	Password string `json:"password"`
}

func (s *Storage) InsertUsers(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudUsers(insertFuncUsers, jsonRequest)
	return id, err

}

func (s *Storage) UpdateUsers(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudUsers(updateFuncUsers, jsonRequest)
	return id, err

}

func (s *Storage) DeleteUsers(jsonRequest map[string]interface{}) (int, error) {

	id, err := s.iudUsers(deleteFuncUsers, jsonRequest)
	return id, err

}

func (s *Storage) iudUsers(nameFunction string, jsonRequest map[string]interface{}) (int, error) {

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

func (s *Storage) ViewUsers(jsonRequest map[string]interface{}) ([]Users, error) {

	rows, err := s.db.Query(context.Background(), "SELECT * FROM "+viewFuncUsers+"($1);", jsonRequest)
	if err != nil {
		return nil, err
	}

	var users []Users

	for rows.Next() {
		var t Users
		err = rows.Scan(
			&t.Id,
			&t.Name,
			&t.Mail,
			&t.Telega,
			&t.Password,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, t)
	}
	return users, rows.Err()
}
