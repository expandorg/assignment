package assignment

type Settings struct {
	ID        uint64 `json:"id" db:"id"`
	JobID     uint64 `json:"job_id" db:"job_id"`       // unique, only one setting per job
	Limit     int    `json:"limit" db:"limit"`         // total assignment limit per job per worker
	Repeat    bool   `json:"repeat" db:"repeat"`       // can a worker repeat the job
	Singly    bool   `json:"singly" db:"singly"`       // single task per job
	Whitelist bool   `json:"whitelist" db:"whitelist"` // job has a whitelist for workers
}
