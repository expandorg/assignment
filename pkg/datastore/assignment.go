package datastore

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/whitelist"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Storage interface {
	GetAssignments(assignment.Params) (assignment.Assignments, error)
	GetAssignment(id string) (*assignment.Assignment, error)
	CreateAssignment(assignment.NewAssignment) (*assignment.Assignment, error)
	GetSettings(jobID uint64) (*assignment.Settings, error)
	GetWhitelist(jobID uint64, workerID uint64) (*whitelist.Whitelist, error)
	WorkerAlreadyAssigned(jobID uint64, workerID uint64) (bool, error)
	DeleteAssignment(workerID uint64, jobID uint64) (bool, error)
	DeactivateAssignment(workerID uint64, jobID uint64) (bool, error)
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
	result, err := as.DB.Exec(
		"INSERT INTO assignments (job_id, task_id, worker_id) VALUES (?,?,?)",
		a.JobID, a.TaskID, a.WorkerID)

	if err != nil {
		if err != nil {
			mysqlerr, ok := err.(*mysql.MySQLError)
			// duplicate entry worker_id & job_id
			if ok && mysqlerr.Number == 1062 {
				return nil, AlreadyAssigned{}
			}
		}
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	assi, err := as.GetAssignment(strconv.FormatInt(id, 10))

	if err != nil {
		return nil, err
	}

	return assi, nil
}

func (as *AssignmentStore) GetSettings(jobID uint64) (*assignment.Settings, error) {
	set := &assignment.Settings{}
	err := as.DB.Get(set, "SELECT * FROM settings WHERE job_id = ?", jobID)

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (as *AssignmentStore) GetWhitelist(jobID uint64, workerID uint64) (*whitelist.Whitelist, error) {
	wl := &whitelist.Whitelist{}
	err := as.DB.Get(wl, "SELECT * FROM whitelists WHERE job_id = ? AND worker_id = ?", jobID, workerID)

	if err != nil {
		return nil, err
	}

	return wl, nil
}

func (as *AssignmentStore) WorkerAlreadyAssigned(jobID uint64, workerID uint64) (bool, error) {
	a := &assignment.Assignment{}
	err := as.DB.Get(a, "SELECT * FROM assignments WHERE job_id = ? AND worker_id = ? AND active IS TRUE", jobID, workerID)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (as *AssignmentStore) DeleteAssignment(workerID uint64, jobID uint64) (bool, error) {
	result, err := as.DB.Exec("DELETE FROM assignments WHERE worker_id = ? AND job_id = ?", workerID, jobID)
	if err != nil {
		return false, err
	}

	numAffected, err := result.RowsAffected()

	if numAffected == 0 {
		return false, AssignmentNotFound{workerID, jobID}
	}

	return true, nil
}

func (as *AssignmentStore) DeactivateAssignment(workerID uint64, jobID uint64) (bool, error) {
	result, err := as.DB.Exec("UPDATE assignments SET active = FALSE WHERE worker_id = ? AND job_id = ?", workerID, jobID)
	if err != nil {
		return false, err
	}

	numAffected, err := result.RowsAffected()

	if numAffected == 0 {
		return false, AssignmentNotFound{workerID, jobID}
	}

	return true, nil
}
