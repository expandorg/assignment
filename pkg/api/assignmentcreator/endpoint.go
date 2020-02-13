package assignmentcreator

import (
	"context"
	"fmt"

	"github.com/expandorg/assignment/pkg/apierror"
	"github.com/expandorg/assignment/pkg/assignment"
	"github.com/expandorg/assignment/pkg/authentication"
	ds "github.com/expandorg/assignment/pkg/datastore"
	"github.com/expandorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AssignmentCreationResponse struct {
	Current     *assignment.Assignment
	Assignments assignment.Assignments
}

func makeAssignmentCreatorEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(assignment.NewAssignment)
		settings, err := svc.GetSettings(req.JobID)
		if err != nil {
			if _, ok := err.(ds.NoRowErr); !ok {
				return nil, errorResponse(err)
			}
		}

		current, err := svc.CreateAssignment(req, settings)
		if err != nil {
			return nil, errorResponse(err)
		}
		params := assignment.Params{
			WorkerID: fmt.Sprintf("%d", req.WorkerID),
			Status:   assignment.Active,
		}
		assignments, err := svc.GetAssignments(params)
		if err != nil {
			return nil, errorResponse(err)
		}

		return AssignmentCreationResponse{Current: current, Assignments: assignments}, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
