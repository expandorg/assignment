package datastore

import (
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetAssignments() (assignment.Assignments, error)
}

type AssignmentStore struct {
	DB *sqlx.DB
}

func NewAssignmentStore(db *sqlx.DB) *AssignmentStore {
	return &AssignmentStore{
		DB: db,
	}
}

func (as *AssignmentStore) GetAssignments() (assignment.Assignments, error) {
	assignments := assignment.Assignments{}

	err := as.DB.Select(&assignments, "SELECT * FROM assignments")
	if err != nil {
		return assignments, err
	}

	return assignments, nil
}
