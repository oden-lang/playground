package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

type Event struct {
	Message string `json:"Message"`
	Kind    string `json:"Kind"`
	Delay   int    `json:"Delay"`
}
type PlayResponse struct {
	Errors string  `json:"Errors"`
	Events []Event `json:"Events"`
}

const compileUrl = "https://play.golang.org/compile"

func runGoPkg(code string, odenVersion string) (*PlayResponse, error) {
	client := &http.Client{}
	form := url.Values{
		"version": {"2"},
		"body":    {code},
	}
	encodedForm := form.Encode()
	req, err := http.NewRequest("POST", compileUrl, bytes.NewBufferString(encodedForm))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Content-Length", strconv.Itoa(len([]byte(encodedForm))))
	req.Header.Add("User-Agent", fmt.Sprintf("oden-playground/%s", odenVersion))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == 200 {
		var response PlayResponse
		err = json.NewDecoder(resp.Body).Decode(&response)
		if err != nil {
			return nil, err
		} else {
			return &response, nil
		}
	}

	if body, err := ioutil.ReadAll(resp.Body); err != nil {
		return nil, errors.New("Failed to run Go code: " + string(body))
	} else {
		return nil, errors.New("Failed to run Go code, status: " + resp.Status)
	}
}
