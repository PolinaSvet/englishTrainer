package storage

import (
	"dictionary/pkg/storage/postgres"
)

type Interface interface {
	InsertUser(postgres.User) (int, error)
	SelectUsers(int) ([]postgres.User, error)
	Close()
}
