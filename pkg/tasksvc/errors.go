package tasksvc

import "fmt"

type NoAvailableTasks struct {
	JobID uint64
}

func (err NoAvailableTasks) Error() string {
	return fmt.Sprintf("No available tasks for job: %d", err.JobID)
}
