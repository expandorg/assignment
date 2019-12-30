package assignment

type JobLimitReached struct{}

func (err JobLimitReached) Error() string {
	return "We reached the limit of total job assignments"
}

type NoAssignmentRepeat struct{}

func (err NoAssignmentRepeat) Error() string {
	return "Multiple assignments are not allowed"
}

type NoResponseRepeat struct{}

func (err NoResponseRepeat) Error() string {
	return "Worker already responded to job"
}

type OnboardingFailure struct{}

func (err OnboardingFailure) Error() string {
	return "Onboarding is either in progress or failed"
}

type WorkerNotEnoughFunds struct{}

func (err WorkerNotEnoughFunds) Error() string {
	return "Worker doesn't have enough funds"
}
