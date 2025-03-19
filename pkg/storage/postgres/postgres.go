package postgres

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/jackc/pgx/v4/pgxpool"
)

// Data store
type Storage struct {
	db *pgxpool.Pool
}

// The constructor accepts a connection string to the database
func New(constr string) (*Storage, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Storage{
		db: db,
	}
	return &s, nil
}

func (s *Storage) Close() {
	s.db.Close()
}

func (s *Storage) ExecuteSQLFromFile(filepath string) error {
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	_, err = s.db.Exec(context.Background(), string(data))
	if err != nil {
		return fmt.Errorf("failed to execute SQL: %w", err)
	}

	return nil
}

type SqlResponse struct {
	ID  int    `json:"id"`
	Err string `json:"err"`
}
