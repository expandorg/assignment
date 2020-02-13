package tasksvc

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/expandorg/assignment/pkg/apierror"
)

type SearchParams struct {
	Random bool
}

type SearchResult struct {
	ID uint64 `json:"id"`
}

func Search(jobID uint64, authToken string, params SearchParams) (SearchResult, error) {
	r := SearchResult{}
	p := fmt.Sprintf("random=%t", params.Random)
	route := fmt.Sprintf("jobs/%d/tasks/search?%s", jobID, p)

	result, err := serviceRequest("GET", route, authToken, nil)
	if err != nil {
		return r, err
	}

	err = json.Unmarshal(result, &r)
	if err != nil {
		return r, err
	}

	if r.ID == 0 {
		return r, NoAvailableTasks{JobID: jobID}
	}

	return r, nil
}

func ValidateTask(taskID uint64, authToken string) error {
	route := fmt.Sprintf("tasks/%d", taskID)

	_, err := serviceRequest("GET", route, authToken, nil)
	if err != nil {
		return err
	}

	return nil
}

func serviceRequest(action, route, authToken string, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}
	serviceURL := fmt.Sprintf("%s/%s", os.Getenv("TASK_SVC_URL"), route)

	req, err := http.NewRequest(action, serviceURL, reqBody)
	if err != nil {
		return nil, errorResponse(err)
	}

	req.AddCookie(&http.Cookie{Name: "JWT", Value: authToken})

	r, err := client.Do(req)
	if err != nil {
		return nil, errorResponse(err)
	}

	if r.StatusCode != 200 {
		return nil, errors.New("bad request")
	}

	body, err := ioutil.ReadAll(r.Body)
	// decoder := json.NewDecoder(r.Body)
	// err = decoder.Decode(&body)
	if err != nil {
		return nil, errorResponse(err)
	}
	return body, nil
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
