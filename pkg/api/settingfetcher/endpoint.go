package settingfetcher

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/expandorg/assignment/pkg/apierror"
	"github.com/expandorg/assignment/pkg/authentication"
	ds "github.com/expandorg/assignment/pkg/datastore"
	"github.com/expandorg/assignment/pkg/service"
	"github.com/go-kit/kit/endpoint"
)

func makeSettingFetcherEndpoint(svc service.AssignmentService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		data, _ := authentication.ParseAuthData(ctx)
		svc.SetAuthData(data)
		req := request.(SettingRequest)

		jobID, err := strconv.ParseUint(req.JobID, 10, 64)
		if err != nil {
			return nil, errorResponse(err)
		}

		s, err := svc.GetSettings(jobID)
		if _, ok := err.(ds.NoRowErr); ok {
			return json.RawMessage("{}"), nil
		}
		if err != nil {
			return nil, errorResponse(err)
		}
		return s, nil
	}
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}

type SettingRequest struct {
	JobID string
}
