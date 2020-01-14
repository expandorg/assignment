package settingfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gemsorg/assignment/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeSettingFetcherEndpoint(s),
		decodeSettingFetcherRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeSettingFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var ok bool
	jobID, ok := vars["job_id"]
	if !ok {
		return nil, fmt.Errorf("missing job_id parameter")
	}
	return SettingRequest{JobID: jobID}, nil
}
