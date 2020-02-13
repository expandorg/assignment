package assignmentdestroyer

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/expandorg/assignment/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentDestroyerEndpoint(s),
		decodeAssignmentDestroyerRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAssignmentDestroyerRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var ok bool
	assignmentID, ok := vars["assignment_id"]
	if !ok {
		return nil, fmt.Errorf("missing assignment_id parameter")
	}
	return AssignmentRequest{AssignmentID: assignmentID}, nil
}
