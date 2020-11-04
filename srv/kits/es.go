package kits

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	domain string = "http://localhost:9200"
)

func ESReqWithJSON(jsonBody []byte, method, endpoint string) ([]byte, error) {

	endpoint = domain + endpoint

	// Sample JSON document to be included as the request body
	body := strings.NewReader(string(jsonBody))

	// An HTTP client for sending the request
	client := &http.Client{}

	// Form the HTTP request
	req, err := http.NewRequest(method, endpoint, body)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Print("es request status: " + resp.Status + "\n")

	return ioutil.ReadAll(resp.Body)
}

func ESReq(method, endpoint string) ([]byte, error) {
	endpoint = domain + endpoint

	client := &http.Client{}

	req, err := http.NewRequest(method, endpoint, nil)

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	fmt.Print("es request status: " + resp.Status + "\n")

	return ioutil.ReadAll(resp.Body)
}
