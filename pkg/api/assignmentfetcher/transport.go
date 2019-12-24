package assignmentfetcher

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
)

func MakeAssignmentFetcherHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentsFetcherEndpoint(s),
		decodeAssignmentsFetcherRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAssignmentsFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	return assignment.Assignments{}, nil
}

// func decodeAssignmentsFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
// 	dr := WorkerDisputesRequest{}
// 	vars := mux.Vars(r)
// 	workerID, ok := vars["worker_id"]
// 	if !ok {
// 		return nil, errorResponse(&apierror.ErrBadRouting{Param: "dispute_id"})
// 	}
// 	id, err := strconv.ParseUint(workerID, 10, 64)
// 	if err != nil {
// 		return nil, errorResponse(&apierror.ErrBadRouting{Param: "dispute_id"})
// 	}
// 	dr.WorkerID = id
// 	return dr, nil
// }
