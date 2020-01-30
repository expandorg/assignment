package service

type NoAvailableTasks struct{}

func (err NoAvailableTasks) Error() string {
	return "There are no available tasks"
}
