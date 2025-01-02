package powerdnsclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// PowerDNSClient stores connection required to maintain connections with PoweDNS API
type PowerDNSClient struct {
	APIConnectionString string
	PowerDNSAPIKey      string
}

// NewClient creates a new PowerDNSClient instance, conection to API is tested
func NewClient(client http.Client, addr string, port int, apikey string) (PowerDNSClient, error) {
	var updater PowerDNSClient
	updater.PowerDNSAPIKey = apikey

	updater.APIConnectionString = fmt.Sprintf("http://%s:%d", addr, port)

	// Checks API connection
	checkURL := fmt.Sprintf("%s/api", updater.APIConnectionString)

	req, err := http.NewRequest(http.MethodGet, checkURL, nil)
	if err != nil {
		return updater, err
	}

	req.Header.Set("X-API-Key", updater.PowerDNSAPIKey)

	response, getErr := client.Do(req)
	if getErr != nil {
		return updater, getErr
	}

	if response.StatusCode != 200 {
		errorString := fmt.Sprintf("Error connecting to API: %d", response.StatusCode)
		return updater, errors.New(errorString)
	}

	return updater, nil
}

// UpdateRecord Updates required DNS records, only A types are supported
func (updater *PowerDNSClient) UpdateRecord(client http.Client, zone string, record string, value string) error {

	updateURL := fmt.Sprintf("%s/api/v1/servers/localhost/zones/%s.", updater.APIConnectionString, zone)
	canonicalRecord := fmt.Sprintf("%s.", record)
	// Request body
	data := map[string]interface{}{
		"rrsets": []map[string]interface{}{
			{
				"name":       canonicalRecord,
				"type":       "A",
				"ttl":        60,
				"changetype": "REPLACE",
				"records": []map[string]interface{}{
					{
						"content":  value,
						"disabled": false,
					},
				},
			},
		},
	}

	jsonData, _ := json.Marshal(data)

	// Prepare Request

	req, reqErr := http.NewRequest(http.MethodPatch, updateURL, bytes.NewBuffer(jsonData))
	if reqErr != nil {
		return reqErr
	}

	// Set headers
	req.Header.Set("X-API-Key", updater.PowerDNSAPIKey)
	req.Header.Set("Content-Type", "application/json")

	resp, respErr := client.Do(req)
	if respErr != nil {
		return respErr
	}

	if resp.StatusCode != http.StatusNoContent {
		return errors.New("Update request failed")
	}

	return nil
}
