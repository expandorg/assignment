package assignmentupdater

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/expandorg/assignment/pkg/apierror"
	"github.com/expandorg/assignment/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentUpdaterEndpoint(s),
		decodeAssignmentUpdaterRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAssignmentUpdaterRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var a AssignmentRequest
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&a)
	if err != nil {
		return nil, apierror.New(500, err.Error(), err)
	}
	return a, nil
}
