package service

import (
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/authorization"
	"github.com/gemsorg/assignment/pkg/datastore"
	"github.com/gemsorg/assignment/pkg/tasksvc"
)

type AssignmentService interface {
	Healthy() bool
	SetAuthData(data authentication.AuthData)
	GetAssignments(assignment.Params) (assignment.Assignments, error)
	GetAssignment(id string) (*assignment.Assignment, error)
	CreateAssignment(assignment.NewAssignment, *assignment.Settings) (*assignment.Assignment, error)
	GetSettings(jobID uint64) (*assignment.Settings, error)
	DeleteAssignment(workerID uint64, jobID uint64) (bool, error)
	DeactivateAssignment(workerID uint64, jobID uint64) (bool, error)
}

type service struct {
	store      datastore.Storage
	authorizor authorization.Authorizer
}

func New(s datastore.Storage, a authorization.Authorizer) *service {
	return &service{
		store:      s,
		authorizor: a,
	}
}

func (s *service) Healthy() bool {
	return true
}

func (s *service) SetAuthData(data authentication.AuthData) {
	s.authorizor.SetAuthData(data)
}

func (s *service) GetAssignments(p assignment.Params) (assignment.Assignments, error) {
	return s.store.GetAssignments(p)
}

func (s *service) GetAssignment(id string) (*assignment.Assignment, error) {
	return s.store.GetAssignment(id)
}

func (s *service) CreateAssignment(a assignment.NewAssignment, set *assignment.Settings) (*assignment.Assignment, error) {
	if set == nil {
		return s.store.CreateAssignment(a)
	}
	// if job has a whitelist, check if worker is part of it
	if set.Whitelist {
		wl, err := s.store.GetWhitelist(a.JobID, a.WorkerID)
		if wl == nil || err != nil {
			return nil, assignment.WorkerNotWhitelisted{}
		}
	}

	// Get worker's assignment for this job
	assigned, err := s.store.WorkerAlreadyAssigned(a.JobID, a.WorkerID)
	a.WorkerAlreadyAssigned = assigned

	// Check if the assignment is allowed
	allowed, err := a.IsAllowed(set)
	if !allowed {
		return nil, err
	}

	// Get a task from the task service
	params := tasksvc.SearchParams{
		Random: true,
	}
	result, err := tasksvc.Search(a.JobID, s.authorizor.GetAuthToken(), params)

	if err != nil {
		return nil, err
	}
	a.TaskID = result.ID

	return s.store.CreateAssignment(a)
}

func (s *service) GetSettings(jobID uint64) (*assignment.Settings, error) {
	return s.store.GetSettings(jobID)
}

func (s *service) DeleteAssignment(workerID uint64, jobID uint64) (bool, error) {
	return s.store.DeleteAssignment(workerID, jobID)
}

func (s *service) DeactivateAssignment(workerID uint64, jobID uint64) (bool, error) {
	return s.store.DeactivateAssignment(workerID, jobID)
}
