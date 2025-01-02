//go:build integration_tests || unit_tests

package powerdnsclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

type RoundTripperMock struct {
	Response *http.Response
	RespErr  error
}

func (rtm *RoundTripperMock) RoundTrip(*http.Request) (*http.Response, error) {
	return rtm.Response, rtm.RespErr
}

func TestInvalidServerConnection(t *testing.T) {

	httpClient := http.Client{
		Timeout: time.Second * 1, // Maximum of 1 second
	}

	var addr string = "327.0.0.1"
	var port int = 8080
	var apiKey = "anykey"

	_, err := NewClient(httpClient, addr, port, apiKey)

	if err == nil {
		t.Errorf("TestInvalidServerConnection should fail.")
	}

}

func TestErroredConnection(t *testing.T) {
	httpClient := http.Client{Transport: &RoundTripperMock{Response: &http.Response{Body: ioutil.NopCloser(bytes.NewBufferString(`Any responde
		`))}, RespErr: errors.New("Error")}}

	var addr string = "127.0.0.1"
	var port int = 8080
	var apiKey = "anykey"

	_, err := NewClient(httpClient, addr, port, apiKey)

	if err == nil {
		t.Errorf("TestInvalidServerConnection should fail.")
	}

}

func TestNonOKResponse(t *testing.T) {

	response := http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Body:       ioutil.NopCloser(bytes.NewBufferString("Unauthorized")),
		Header:     make(http.Header),
	}

	httpClient := http.Client{Transport: &RoundTripperMock{Response: &response}}

	var addr string = "127.0.0.1"
	var port int = 8080
	var apiKey = "anykey"

	_, err := NewClient(httpClient, addr, port, apiKey)

	if err == nil {
		t.Errorf("TestInvalidServerConnection should fail.")
	}

}

func TestOKResponse(t *testing.T) {

	data := map[string]interface{}{
		"message": "Success",
		"status":  "ok",
	}
	jsonBody, _ := json.Marshal(data) // Convert map to JSON string

	// Create an http.Response instance
	response := http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Body:       ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
		Header:     make(http.Header),
	}

	// Add headers (optional)
	response.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{Transport: &RoundTripperMock{Response: &response}}

	var addr string = "127.0.0.1"
	var port int = 8080
	var apiKey = "anykey"

	_, err := NewClient(httpClient, addr, port, apiKey)

	if err != nil {
		t.Errorf("TestInvalidServerConnection should not fail.")
	}

}

func TestFailedUpdate(t *testing.T) {
	data := map[string]interface{}{
		"message": "Fail",
	}
	jsonBody, _ := json.Marshal(data) // Convert map to JSON string

	// Create an http.Response instance
	response := http.Response{
		Status:     "401 Unauthorized",
		StatusCode: http.StatusUnauthorized,
		Body:       ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
		Header:     make(http.Header),
	}

	response.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{Transport: &RoundTripperMock{Response: &response}}

	clientInstance := PowerDNSClient{APIConnectionString: "http:127.0.0.1:8080", PowerDNSAPIKey: "anykey"}

	err := clientInstance.UpdateRecord(httpClient, "windmaker.net", "ejemplo.windmaker.net", "128.34.21.21")

	if err == nil {
		t.Errorf("testFailedUpdate should fail.")
	}

}

func TestSuccesfulUpdate(t *testing.T) {
	data := map[string]interface{}{
		"message": "Fail",
	}
	jsonBody, _ := json.Marshal(data) // Convert map to JSON string

	// Create an http.Response instance
	response := http.Response{
		Status:     "204 No Content",
		StatusCode: http.StatusNoContent,
		Body:       ioutil.NopCloser(bytes.NewBuffer(jsonBody)),
		Header:     make(http.Header),
	}

	response.Header.Set("Content-Type", "application/json")

	httpClient := http.Client{Transport: &RoundTripperMock{Response: &response}}

	clientInstance := PowerDNSClient{APIConnectionString: "http:127.0.0.1:8080", PowerDNSAPIKey: "anykey"}

	err := clientInstance.UpdateRecord(httpClient, "windmaker.net", "ejemplo.windmaker.net", "128.34.21.21")

	if err != nil {
		t.Errorf("testFailedUpdate should not fail.")
	}

}
