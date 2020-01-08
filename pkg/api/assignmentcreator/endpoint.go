package assignmentcreator

import (
	"context"

	"github.com/gemsorg/assignment/pkg/apierror"
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/authentication"
	ds "github.com/gemsorg/assignment/pkg/datastore"
	"github.com/gemsorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

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

		saved, err := svc.CreateAssignment(req, settings)
		if err != nil {
			return nil, errorResponse(err)
		}
		return saved, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
