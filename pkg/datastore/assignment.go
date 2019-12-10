package datastore

import (
	"github.com/jmoiron/sqlx"
)

type Storage interface {
}

type AssignmentStore struct {
	DB *sqlx.DB
}

func NewAssignmentStore(db *sqlx.DB) *AssignmentStore {
	return &AssignmentStore{
		DB: db,
	}
}
