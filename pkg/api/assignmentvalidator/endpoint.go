package assignmentvalidator

import (
	"context"

	"github.com/expandorg/assignment/pkg/apierror"
	"github.com/expandorg/assignment/pkg/assignment"
	"github.com/expandorg/assignment/pkg/authentication"
	ds "github.com/expandorg/assignment/pkg/datastore"
	"github.com/expandorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AssignmentValidatorResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

func makeAssignmentValidatorEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		res := AssignmentValidatorResponse{}
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)

		req := request.(assignment.NewAssignment)
		settings, err := svc.GetSettings(req.JobID)
		if err != nil {
			if _, ok := err.(ds.NoRowErr); !ok {
				return nil, errorResponse(err)
			}
		}

		valid, err := svc.ValidateAssignment(req, settings)
		if err != nil {
			res.Reason = err.Error()
		}
		res.Valid = valid

		return res, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
