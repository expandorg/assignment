package datastore

import "fmt"

type AlreadyAssigned struct{}

func (err AlreadyAssigned) Error() string {
	return "User is already assigned a task from this job"
}

type AlreadyHasSettings struct{}

func (err AlreadyHasSettings) Error() string {
	return "Job already has settings"
}

type AssignmentNotFound struct {
	ID         string
	WorkerID   uint64
	JobID      uint64
	ResponseID uint64
}

func (err AssignmentNotFound) Error() string {
	return fmt.Sprintf("No Record found for id: %s, worker_id: %d, job_id: %d, response_id: %d", err.ID, err.WorkerID, err.JobID, err.ResponseID)
}

type NoRowErr struct{}

func (err NoRowErr) Error() string {
	return "No Records"
}

type RecordsMismatch struct{}

func (err RecordsMismatch) Error() string {
	return "Some records where not processed properly"
}
