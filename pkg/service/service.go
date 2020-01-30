package service

import (
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/authorization"
	"github.com/gemsorg/assignment/pkg/datastore"
	"github.com/gemsorg/assignment/pkg/externalsvc"
	"github.com/gemsorg/assignment/pkg/registrysvc"
	"github.com/gemsorg/assignment/pkg/tasksvc"
)

type AssignmentService interface {
	Healthy() bool
	SetAuthData(data authentication.AuthData)
	GetAssignments(assignment.Params) (assignment.Assignments, error)
	GetAssignment(id string) (*assignment.Assignment, error)
	CreateAssignment(assignment.NewAssignment, *assignment.Settings) (*assignment.Assignment, error)
	GetSettings(jobID uint64) (*assignment.Settings, error)
	DeleteAssignment(id string) (bool, error)
	DeleteAssignments(ids []string) error
	UpdateAssignment(workerID, jobID, responseID uint64, status string) (bool, error)
	CreateSettings(assignment.Settings) (*assignment.Settings, error)
	GetStore() datastore.Storage
	ValidateAssignment(a assignment.NewAssignment, set *assignment.Settings) (bool, error)
	GetRegistration(jobID uint64) (*registrysvc.Registration, error)
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

func (s *service) GetStore() datastore.Storage {
	return s.store
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
	// Get worker's assignment for this job
	assigned, err := s.store.WorkerAlreadyAssigned(a.JobID, a.WorkerID)
	a.WorkerAlreadyAssigned = assigned

	// Check if there's an external service registered for this task
	r, err := s.GetRegistration(a.JobID)
	if r != nil && r.Services[registrysvc.AssignmentCreator] != nil {
		res, err := externalsvc.Assign(*r, a, s.authorizor.GetAuthToken())
		if res == nil || err != nil {
			return nil, NoAvailableTasks{}
		}
		a.TaskID = res.TaskID
		return s.store.CreateAssignment(a)
	}

	if set != nil {
		// if job has a whitelist, check if worker is part of it
		if set.Whitelist {
			wl, err := s.store.GetWhitelist(a.JobID, a.WorkerID)
			if wl == nil || err != nil {
				return nil, assignment.WorkerNotWhitelisted{}
			}
		}

		// Check if the assignment is allowed
		allowed, err := a.IsAllowed(set)
		if !allowed {
			return nil, err
		}
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

func (s *service) DeleteAssignment(id string) (bool, error) {
	return s.store.DeleteAssignment(id)
}

func (s *service) DeleteAssignments(ids []string) error {
	return s.store.DeleteAssignments(ids)
}

func (s *service) UpdateAssignment(workerID, jobID, responseID uint64, status string) (bool, error) {
	return s.store.UpdateAssignment(workerID, jobID, responseID, status)
}

func (s *service) CreateSettings(set assignment.Settings) (*assignment.Settings, error) {
	return s.store.CreateSettings(set)
}

func (s *service) ValidateAssignment(a assignment.NewAssignment, set *assignment.Settings) (bool, error) {
	allowed := true
	if set == nil {
		return allowed, nil
	}

	// if job has a whitelist, check if worker is part of it
	if set.Whitelist {
		wl, err := s.store.GetWhitelist(a.JobID, a.WorkerID)
		if wl == nil || err != nil {
			return false, assignment.WorkerNotWhitelisted{}
		}
	}

	// Get worker's assignment for this job
	assigned, err := s.store.WorkerAlreadyAssigned(a.JobID, a.WorkerID)
	a.WorkerAlreadyAssigned = assigned

	// Check if there's an external service registered for this task
	r, err := s.GetRegistration(a.JobID)
	if r != nil && r.Services[registrysvc.AssignmentValidator] != nil {
		return externalsvc.Validate(*r, a)
	} else {
		allowed, err = a.IsAllowed(set)
	}

	// Check if the assignment is allowed
	if !allowed {
		return false, err
	}

	return allowed, nil
}

func (s *service) GetRegistration(jobID uint64) (*registrysvc.Registration, error) {
	return registrysvc.GetRegistration(s.authorizor.GetAuthToken(), jobID)
}
