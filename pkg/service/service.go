package service

import (
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/authorization"
	"github.com/gemsorg/assignment/pkg/datastore"
)

type AssignmentService interface {
	Healthy() bool
	SetAuthData(data authentication.AuthData)
	GetAssignments(assignment.Params) (assignment.Assignments, error)
	GetAssignment(id string) (*assignment.Assignment, error)
	CreateAssignment(assignment.NewAssignment, *assignment.Settings) (*assignment.Assignment, error)
	GetSettings(jobID uint64) (*assignment.Settings, error)
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

	return s.store.CreateAssignment(a)
}

func (s *service) GetSettings(jobID uint64) (*assignment.Settings, error) {
	return s.store.GetSettings(jobID)
}
