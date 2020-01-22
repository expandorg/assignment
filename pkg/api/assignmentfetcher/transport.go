package assignmentfetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/service"
	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
)

func MakeAssignmentsFetcherHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentsFetcherEndpoint(s),
		decodeAssignmentsFetcherRequest,
		encodeResponse,
	)
}

func MakeAssignmentFetcherHandler(s service.AssignmentService) http.Handler {
	return kithttp.NewServer(
		makeAssignmentFetcherEndpoint(s),
		decodeAssignmentFetcherRequest,
		encodeResponse,
	)
}

func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

func decodeAssignmentsFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	as := assignment.Params{}
	params := r.URL.Query()

	workerID, ok := params["worker_id"]
	if ok && len(workerID) > 0 {
		as.WorkerID = workerID[0]
	}
	jobID, ok := params["job_id"]
	if ok && len(jobID) > 0 {
		as.JobID = jobID[0]
	}
	taskID, ok := params["task_id"]
	if ok && len(taskID) > 0 {
		as.TaskID = taskID[0]
	}
	responseID, ok := params["response_id"]
	if ok && len(responseID) > 0 {
		as.ResponseID = responseID[0]
	}
	status, ok := params["status"]
	if ok && len(taskID) > 0 {
		as.Status = assignment.Status(status[0])
	}

	return as, nil
}

func decodeAssignmentFetcherRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	var ok bool
	assignmentID, ok := vars["assignment_id"]
	if !ok {
		return nil, fmt.Errorf("missing assignment_id parameter")
	}
	return AssignmentRequest{AssignmentID: assignmentID}, nil
}
