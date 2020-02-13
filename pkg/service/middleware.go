package service

import (
	"context"
	"time"

	"github.com/go-kit/kit/log"

	assignment "github.com/expandorg/assignment/pkg/assignment"
	authentication "github.com/expandorg/assignment/pkg/authentication"
	datastore "github.com/expandorg/assignment/pkg/datastore"
	registrysvc "github.com/expandorg/assignment/pkg/registrysvc"
)

type Middleware func(AssignmentService) AssignmentService

type loggingMiddleware struct {
	next   AssignmentService
	logger log.Logger
}

func ServiceMiddleware(logger log.Logger) Middleware {
	return func(next AssignmentService) AssignmentService {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func (mw loggingMiddleware) CreateAssignment(ctx context.Context, asgn assignment.Asgn, set *assignment.Settings) (*assignment.Assignment, error) {
	defer func(begin time.Time) {
		mw.logger.Log("method", "CreateAssignment", "id", "took", time.Since(begin))
	}(time.Now())
	return mw.next.CreateAssignment(ctx, asgn, set)
}

func (mw loggingMiddleware) Healthy() bool {
	return mw.next.Healthy()
}

func (mw loggingMiddleware) SetAuthData(data authentication.AuthData) {
	mw.next.SetAuthData(data)
}
func (mw loggingMiddleware) GetAssignments(a assignment.Params) (assignment.Assignments, error) {
	return mw.next.GetAssignments(a)
}
func (mw loggingMiddleware) GetAssignment(id string) (*assignment.Assignment, error) {
	return mw.next.GetAssignment(id)
}
func (mw loggingMiddleware) GetSettings(jobID uint64) (*assignment.Settings, error) {
	return mw.next.GetSettings(jobID)
}
func (mw loggingMiddleware) DeleteAssignment(id string) (bool, error) {
	return mw.next.DeleteAssignment(id)
}
func (mw loggingMiddleware) DeleteAssignments(ids []string) error {
	return mw.next.DeleteAssignments(ids)
}
func (mw loggingMiddleware) UpdateAssignment(workerID, jobID, responseID uint64, status string) (bool, error) {
	return mw.next.UpdateAssignment(workerID, jobID, responseID, status)
}
func (mw loggingMiddleware) CreateSettings(s assignment.Settings) (*assignment.Settings, error) {
	return mw.next.CreateSettings(s)
}
func (mw loggingMiddleware) GetStore() datastore.Storage {
	return mw.next.GetStore()
}
func (mw loggingMiddleware) ValidateAssignment(a assignment.NewAssignment, set *assignment.Settings) (bool, error) {
	return mw.next.ValidateAssignment(a, set)
}
func (mw loggingMiddleware) GetRegistration(jobID uint64) (*registrysvc.Registration, error) {
	return mw.next.GetRegistration(jobID)
}
