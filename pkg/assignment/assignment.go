package assignment

import (
	"time"

	"github.com/gemsorg/assignment/pkg/nulls"
)

type Params struct {
	WorkerID   string
	JobID      string
	TaskID     string
	ResponseID string
}

type Assignment struct {
	ID         uint64      `db:"id" json:"id"`
	JobID      uint64      `db:"job_id" json:"job_id"`
	TaskID     uint64      `db:"task_id" json:"task_id"`
	ResponseID nulls.Int64 `db:"response_id" json:"response_id"`
	WorkerID   uint64      `db:"worker_id" json:"worker_id"`
	Active     bool        `db:"active" json:"active"`
	AssignedAt time.Time   `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time  `db:"expires_at" json:"expires_at"`
}

type NewAssignment struct {
	JobID                  uint64 `json:"job_id"`
	TaskID                 uint64 `json:"task_id"`
	WorkerID               uint64 `json:"worker_id"`
	JobAssignmentCount     int    `json:"job_assignment_count"`
	OnboardingStatus       bool   `json:"onboarding_status"`
	WorkerAlreadyAssigned  bool   `json:"worker_already_assigned"`
	WorkerAlreadyResponded bool   `json:"worker_already_responded"`
	WorkerHasFunds         bool   `json:"worker_has_funds"`
}

type Assignments []Assignment

func (a NewAssignment) IsAllowed(set *Settings) (bool, error) {
	// check onboarding
	if !a.OnboardingStatus {
		return false, OnboardingFailure{}
	}
	// We reached the limit of the total assignment for the job
	if set.Limit != 0 && set.Limit == a.JobAssignmentCount {
		return false, JobLimitReached{}
	}

	// Only allow the worker to be once at a time
	if set.Singly && a.WorkerAlreadyAssigned {
		return false, NoAssignmentRepeat{}
	}

	// Worker can only respond once for the the job
	if !set.Repeat && a.WorkerAlreadyResponded {
		return false, NoResponseRepeat{}
	}

	if !a.WorkerHasFunds {
		return false, WorkerNotEnoughFunds{}
	}

	return true, nil
}
