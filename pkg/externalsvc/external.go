package externalsvc

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gemsorg/assignment/pkg/assignment"
	"github.com/gemsorg/assignment/pkg/registrysvc"
)

type ValidationResponse struct {
	Valid  bool   `json:"valid"`
	Reason string `json:"reason"`
}

func Validate(reg registrysvc.Registration, a assignment.NewAssignment) (bool, error) {
	esr := ValidationResponse{}
	url := reg.Services[registrysvc.AssignmentValidator].URL
	res, err := serviceRequest("POST", url, reg.APIKeyID, reg.RequesterID)
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

func serviceRequest(action, url, key string, userID uint64) ([]byte, error) {
	client := &http.Client{}

	req, err := http.NewRequest(action, url, nil)
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
