package assignment

import (
	"time"

	"github.com/gemsorg/assignment/pkg/nulls"
)

type Assignment struct {
	ID         uint64     `db:"id" json:"id"`
	JobID      uint64     `db:"job_id" json:"job_id"`
	TaskID     uint64     `db:"task_id" json:"task_id"`
	ResponseID uint64     `db:"response_id" json:"response_id"`
	WorkerID   uint64     `db:"worker_id" json:"worker_id"`
	AssignedAt time.Time  `db:"assigned_at" json:"assigned_at"`
	ExpiresAt  nulls.Time `db:"expires_at" json:"expires_at"`
}

type Assignments []Assignment
