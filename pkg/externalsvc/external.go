package externalsvc

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/registrysvc"
	"github.com/gemsorg/assignment/pkg/tasksvc"
)

type ValidationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

type CreatorResponse struct {
	TaskID uint64 `json:"task_id"`
}

func Validate(reg registrysvc.Registration, a assignment.NewAssignment) (bool, error) {
	esr := ValidationResponse{}
	url := reg.Services[registrysvc.AssignmentValidator].URL
	requestByte, _ := json.Marshal(a)
	reqBody := bytes.NewReader(requestByte)
	res, err := serviceRequest("POST", url, reg.APIKeyID, reg.RequesterID, reqBody)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal(res, &esr)
	if err != nil {
		return false, err
	}

	if !esr.Valid {
		return false, AssignmentRejectedByJobOwner{}
	}

	return true, nil
}

func Assign(reg registrysvc.Registration, a assignment.NewAssignment, authToken string) (*CreatorResponse, error) {
	cr := CreatorResponse{}
	url := reg.Services[registrysvc.AssignmentCreator].URL
	requestByte, _ := json.Marshal(a)
	reqBody := bytes.NewReader(requestByte)

	res, err := serviceRequest("PUT", url, reg.APIKeyID, reg.RequesterID, reqBody)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &cr)
	if err != nil {
		return nil, err
	}

	if cr.TaskID == 0 {
		return nil, AssignmentRejectedByJobOwner{}
	}

	// validate task
	err = tasksvc.ValidateTask(cr.TaskID, authToken)
	if err != nil {
		return nil, AssignmentRejectedByJobOwner{}
	}

	return &cr, nil
}

func serviceRequest(action, url, key string, userID uint64, reqBody io.Reader) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(action, url, reqBody)
	if err != nil {
		return nil, err
	}
	apiKey, err := GenerateAPIKeyJWT(userID, key)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", apiKey)
	r, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
