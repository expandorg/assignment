package whitelist

type Whitelist struct {
	ID       uint64 `db:"id"`
	JobID    uint64 `db:"job_id"`
	WorkerID uint64 `db:"worker_id"`
}
