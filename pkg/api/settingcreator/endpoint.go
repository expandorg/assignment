package settingcreator

import (
	"context"

	"github.com/gemsorg/assignment/pkg/apierror"
	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/authentication"
	"github.com/gemsorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

type AssignmentCreationResponse struct {
	Current     *assignment.Assignment
	Assignments assignment.Assignments
}

func makeSettingsCreatorEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(assignment.Settings)
		settings, err := svc.CreateSettings(req)
		if err != nil {
			return nil, errorResponse(err)
		}
		return settings, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
