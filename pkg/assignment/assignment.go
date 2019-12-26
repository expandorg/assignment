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
	AssignedAt time.Time   `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time  `db:"expires_at" json:"expires_at"`
}

type NewAssignment struct {
	JobID    uint64 `json:"job_id"`
	TaskID   uint64 `json:"task_id"`
	WorkerID uint64 `json:"worker_id"`
}

type Assignments []Assignment
