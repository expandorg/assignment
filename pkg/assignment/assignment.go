package assignment

import (
	"time"

	"github.com/expandorg/assignment/pkg/nulls"
)

type Status string

const (
	Active   Status = "active"
	InActive Status = "inactive"
	Pending  Status = "pending"
	Accepted Status = "accepted"
	Rejected Status = "rejected"
	Expired  Status = "expired"
)

type Params struct {
	WorkerID   string
	JobID      string
	TaskID     string
	ResponseID string
	Status     Status
}

type Assignment struct {
	ID         uint64      `db:"id" json:"id"`
	JobID      uint64      `db:"job_id" json:"job_id"`
	TaskID     uint64      `db:"task_id" json:"task_id"`
	WorkerID   uint64      `db:"worker_id" json:"worker_id"`
	ResponseID nulls.Int64 `db:"response_id" json:"response_id"`
	Status     string      `db:"status" json:"status"`
	Active     nulls.Bool  `db:"active" json:"active"`
	AssignedAt time.Time   `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time  `db:"expires_at" json:"expires_at"`
}

type NewAssignment struct {
	JobID                  uint64 `json:"job_id"`
	TaskID                 uint64 `json:"task_id"`
	WorkerID               uint64 `json:"worker_id"`
	WorkerAssignmentCount  int    // This comes internally from this service
	OnboardingSuccess      bool   `json:"onboarding_success"`
	WorkerAlreadyResponded bool   `json:"worker_already_responded"`
	WorkerAlreadyAssigned  bool   // This comes internally from this service
}

type Assignments []Assignment

func (a NewAssignment) IsAllowed(set *Settings) (bool, error) {
	// check onboarding
	if !a.OnboardingSuccess {
		return false, OnboardingFailure{}
	}

	// We reached the limit of the total assignment for the job and worker
	if set.Limit != 0 && a.WorkerAssignmentCount >= set.Limit {
		return false, JobLimitReached{}
	}

	// Only allow the worker to be assigned once
	if set.Singly && a.WorkerAlreadyAssigned {
		return false, NoAssignmentRepeat{}
	}

	// Worker can only respond once for the the job
	if !set.Repeat && a.WorkerAlreadyResponded {
		return false, NoResponseRepeat{}
	}

	return true, nil
}
