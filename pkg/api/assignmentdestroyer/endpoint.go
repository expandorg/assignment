package assignmentdestroyer

import (
	"context"

	"github.com/gemsorg/assignment/pkg/apierror"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentDestroyerEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(AssignmentRequest)
		p, err := svc.DeleteAssignment(req.AssignmentID)
		if err != nil {
			return AssignmentResponse{p}, errorResponse(err)
		}
		return AssignmentResponse{p}, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type AssignmentRequest struct {
	AssignmentID string `json:"assignment_id"`
}

type AssignmentResponse struct {
	Destroyed bool `json:"destroyed"`
}
