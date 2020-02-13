package registrysvc

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/expandorg/assignment/pkg/apierror"
)

const (
	AssignmentValidator = "AssignmentValidator"
	AssignmentCreator   = "AssignmentCreator"
)

type Registration struct {
	ID          uint64   `json:"id"`
	JobID       uint64   `json:"job_id"`
	APIKeyID    string   `json:"api_key_id"`
	Services    Services `json:"services"`
	RequesterID uint64   `json:"requester_id"`
}

type Services map[string]*Service

type Service struct {
	URL string `json:"url"`
}

func GetRegistration(authToken string, jobID uint64) (*Registration, error) {
	r := &Registration{}

	url := fmt.Sprintf("registrations/%d", jobID)
	res, err := serviceRequest("GET", url, authToken)
	json.Unmarshal(res, r)

	if err != nil {
		return nil, err
	}
	fmt.Printf("R %+v\n", r)
	return r, nil
}

func serviceRequest(action, route, authToken string) ([]byte, error) {
	client := &http.Client{}
	serviceURL := fmt.Sprintf("%s/%s", os.Getenv("REGISTRY_SVC_URL"), route)

	req, err := http.NewRequest(action, serviceURL, nil)
	if err != nil {
		return nil, errorResponse(err)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)
	r, err := client.Do(req)
	if err != nil {
		return nil, errorResponse(err)
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, errorResponse(err)
	}

	return body, nil
}

func errorResponse(err error) *apierror.APIError {
	return apierror.New(500, err.Error(), err)
}
