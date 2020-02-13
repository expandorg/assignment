package datastore

import (
	"strconv"
	"strings"

	"github.com/expandorg/assignment/pkg/assignment"
	"github.com/expandorg/assignment/pkg/whitelist"
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
	DeleteAssignment(id string) (bool, error)
	DeleteAssignments(ids []string) error
	UpdateAssignment(workerID, jobID, responseID uint64, status string) (bool, error)
	CreateSettings(assignment.Settings) (*assignment.Settings, error)
	SelectExpiredAssignments() (assignment.Assignments, error)
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
	// If we're looking for expired assignments, then return them
	if p.Status == assignment.Expired {
		return as.SelectExpiredAssignments()
	}

	assignments := assignment.Assignments{}
	query := "SELECT * FROM assignments"
	paramsQuery := []string{}
	args := []interface{}{}

	if p.WorkerID != "" && p.WorkerID != "0" {
		args = append(args, p.WorkerID)
		paramsQuery = append(paramsQuery, "worker_id=?")
	}
	if p.JobID != "" && p.JobID != "0" {
		args = append(args, p.JobID)
		paramsQuery = append(paramsQuery, "job_id=?")
	}
	if p.TaskID != "" && p.TaskID != "0" {
		args = append(args, p.TaskID)
		paramsQuery = append(paramsQuery, "task_id=?")
	}
	if p.ResponseID != "" && p.ResponseID != "0" {
		args = append(args, p.ResponseID)
		paramsQuery = append(paramsQuery, "response_id=?")
	}
	if p.Status != "" {
		args = append(args, p.Status)
		paramsQuery = append(paramsQuery, "status=?")
	}

	if len(paramsQuery) > 0 {
		query = query + " Where " + strings.Join(paramsQuery, " AND ")
	}

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
		"INSERT INTO assignments (job_id, task_id, worker_id, active, expires_at) VALUES (?,?,?,?,DATE_ADD(CURRENT_TIMESTAMP, INTERVAL 2 HOUR))",
		a.JobID, a.TaskID, a.WorkerID, 1)

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
	set := []*assignment.Settings{}
	err := as.DB.Select(&set, "SELECT * FROM settings WHERE job_id = ?", jobID)

	if err != nil {
		return nil, err
	}

	if len(set) == 0 {
		return nil, NoRowErr{}
	}

	return set[0], nil
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
	err := as.DB.Get(a, "SELECT * FROM assignments WHERE job_id = ? AND worker_id = ? AND status = ?", jobID, workerID, assignment.Active)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (as *AssignmentStore) DeleteAssignment(id string) (bool, error) {
	result, err := as.DB.Exec("DELETE FROM assignments WHERE id = ?", id)
	if err != nil {
		return false, err
	}

	numAffected, err := result.RowsAffected()

	if numAffected == 0 {
		return false, AssignmentNotFound{ID: id}
	}

	return true, nil
}

func (as *AssignmentStore) UpdateAssignment(workerID, jobID, responseID uint64, status string) (bool, error) {
	// we're making it pending from active
	var numAffected int64
	if assignment.Status(status) == assignment.Pending {
		result, err := as.DB.Exec(
			"UPDATE assignments SET status = ?, active = ?, response_id = ? WHERE worker_id = ? AND job_id = ? AND active IS TRUE",
			status, nil, responseID, workerID, jobID,
		)
		if err != nil {
			return false, err
		}
		numAffected, err = result.RowsAffected()
	} else {
		// we're scoring
		result, err := as.DB.Exec(
			"UPDATE assignments SET status = ? WHERE worker_id = ? AND job_id = ? AND response_id = ? AND active IS NULL",
			status, workerID, jobID, responseID,
		)
		if err != nil {
			return false, err
		}
		numAffected, err = result.RowsAffected()
	}

	if numAffected == 0 {
		return false, AssignmentNotFound{WorkerID: workerID, JobID: jobID, ResponseID: responseID}
	}

	return true, nil
}

func (as *AssignmentStore) CreateSettings(s assignment.Settings) (*assignment.Settings, error) {
	// WE always replace the settings with the incoming
	_, err := as.DB.Exec(
		"REPLACE INTO settings (`limit`, `repeat`, singly, whitelist, job_id) VALUES (?, ?, ?, ?, ?)",
		s.Limit, s.Repeat, s.Singly, s.Whitelist, s.JobID,
	)

	if err != nil {
		if err != nil {
			mysqlerr, ok := err.(*mysql.MySQLError)
			// duplicate entry job_id
			if ok && mysqlerr.Number == 1062 {
				return nil, AlreadyHasSettings{}
			}
		}
		return nil, err
	}

	set, err := as.GetSettings(s.JobID)

	if err != nil {
		return nil, err
	}

	return set, nil
}

func (as *AssignmentStore) SelectExpiredAssignments() (assignment.Assignments, error) {
	assignments := assignment.Assignments{}
	err := as.DB.Select(
		&assignments,
		`SELECT * FROM assignments WHERE active IS TRUE AND expires_at <= NOW()`,
	)
	if err != nil {
		return nil, err
	}
	return assignments, nil
}

func (as *AssignmentStore) DeleteAssignments(ids []string) error {
	query, args, _ := sqlx.In(`DELETE FROM assignments WHERE id IN (?)`, ids)
	query = as.DB.Rebind(query)

	res, err := as.DB.Exec(query, args...)

	if err != nil {
		return err
	}

	numRows, _ := res.RowsAffected()

	if numRows != int64(len(ids)) {
		return RecordsMismatch{}
	}

	return nil
}
