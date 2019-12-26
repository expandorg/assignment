package datastore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetAssignments(assignment.Params) (assignment.Assignments, error)
	GetAssignment(id string) (*assignment.Assignment, error)
	CreateAssignment(assignment.NewAssignment) (*assignment.Assignment, error)
}

type AssignmentStore struct {
	DB *sqlx.DB
}

func NewAssignmentStore(db *sqlx.DB) *AssignmentStore {
	return &AssignmentStore{
		DB: db,
	}
}

func (as *AssignmentStore) GetAssignments(p assignment.Params) (assignment.Assignments, error) {
	assignments := assignment.Assignments{}
	query := "SELECT * FROM assignments"
	paramsQuery := []string{}
	args := []interface{}{}
	if p.WorkerID != "" {
		args = append(args, p.WorkerID)
		paramsQuery = append(paramsQuery, "worker_id=?")
	}
	if p.JobID != "" {
		args = append(args, p.JobID)
		paramsQuery = append(paramsQuery, "job_id=?")
	}
	if p.TaskID != "" {
		args = append(args, p.TaskID)
		paramsQuery = append(paramsQuery, "task_id=?")
	}
	if p.ResponseID != "" {
		args = append(args, p.ResponseID)
		paramsQuery = append(paramsQuery, "response_id=?")
	}

	if len(paramsQuery) > 0 {
		query = query + " Where " + strings.Join(paramsQuery, " AND ")
	}
	fmt.Println("Q", query)
	err := as.DB.Select(&assignments, query, args...)
	if err != nil {
		return assignments, err
	}

	return assignments, nil
}
func (as *AssignmentStore) GetAssignment(id string) (*assignment.Assignment, error) {
	assignment := &assignment.Assignment{}
	err := as.DB.Get(assignment, "SELECT * FROM assignments WHERE id = ?", id)

	if err != nil {
		return nil, err
	}

	return assignment, nil
}

func (as *AssignmentStore) CreateAssignment(a assignment.NewAssignment) (*assignment.Assignment, error) {
	newAs := &assignment.Assignment{}
	result, err := as.DB.Exec(
		"INSERT INTO assignments (job_id, task_id, worker_id) VALUES (?,?,?)",
		a.JobID, a.TaskID, a.WorkerID)

	if err != nil {
		return newAs, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return newAs, err
	}

	assi, err := as.GetAssignment(strconv.FormatInt(id, 10))

	if err != nil {
		return newAs, err
	}

	return assi, nil
}
