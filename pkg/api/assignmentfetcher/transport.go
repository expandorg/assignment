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

	return as, nil
}
