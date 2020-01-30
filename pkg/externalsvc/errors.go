package externalsvc

type AssignmentRejectedByJobOwner struct{}

func (err AssignmentRejectedByJobOwner) Error() string {
	return "Job owner has rejected this assignment through external service"
}
