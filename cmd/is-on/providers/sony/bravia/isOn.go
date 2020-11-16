package bravia

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Sharsie/tv-status-rpio/cmd/is-on/config"
)

type requestData struct {
	Method  string   `json:"method"`
	Version string   `json:"version"`
	Id      int      `json:"id"`
	Params  []string `json:"params"`
}

type responseData struct {
	Result []struct {
		Status string `json:"status"`
	} `json:"result"`
	Id int `json:"id"`
}

const statusEndpoint = "/sony/system"

var httpEndpoint = config.TvHostname + statusEndpoint

func IsOn() (bool, error) {
	log.Println("Getting Sony Bravia TV status.")

	payload := requestData{
		"getPowerStatus",
		"1.0",
		1,
		make([]string, 0),
	}

	requestBody, err := json.Marshal(payload)

	if err != nil {
		return false, errors.New("Could not create a request payload.")
	}

	response, err := http.Post(httpEndpoint, "application/json", bytes.NewBuffer(requestBody))

	if err != nil {
		return false, errors.New("Could not get TV response.")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return false, errors.New("Could not read the TV response.")
	}

	var data responseData

	err = json.Unmarshal(body, &data)

	if err != nil || len(data.Result) < 1 {
		return false, errors.New("Could not decode the TV response.")
	}

	log.Printf("The Sony Bravia TV status is '%s'.", data.Result[0].Status)

	return data.Result[0].Status == "active" || data.Result[0].Status == "activating", nil
}
