package datastore

import "fmt"

type AlreadyAssigned struct{}

func (err AlreadyAssigned) Error() string {
	return "User is already assigned a task from this job"
}

type AssignmentNotFound struct {
	WorkerID uint64
	JobID    uint64
}

func (err AssignmentNotFound) Error() string {
	return fmt.Sprintf("No Record found for worker_id: %d, job_id: %d", err.WorkerID, err.JobID)
}

type NoRowErr struct{}

func (err NoRowErr) Error() string {
	return "No Records"
}
