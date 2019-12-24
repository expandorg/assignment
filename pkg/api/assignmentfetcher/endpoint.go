package assignmentfetcher

import (
	"context"

	"github.com/gemsorg/assignment/pkg/apierror"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeAssignmentsFetcherEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		assignments, err := svc.GetAssignments()
		if err != nil {
			return nil, errorResponse(err)
		}
		return assignments, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
