//go:build integration_tests || unit_tests

package updater

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/a-castellano/home-ip-updater/powerdnsclient"
)

var currentAWSAccessKey string
var currentAWSAccessKeyDefined bool

var currentAWSSecretKey string
var currentAWSSecretKeyDefined bool

var currentAWSZoneId string
var currentAWSZoneIdDefined bool

var currentSubdomain string
var currentSubdomainDefined bool

var currentPowerDNSHost string
var currentPowerDNSHostDefined bool

var currentPowerDNSPort string
var currentPowerDNSPortDefined bool

var currentPowerDNSAPIKey string
var currentPowerDNSAPIKeyDefined bool

func setUp() {

	if envAWSAccessKey, found := os.LookupEnv("AWS_ACCESS_KEY_ID"); found {
		currentAWSAccessKey = envAWSAccessKey
		currentAWSAccessKeyDefined = true
	} else {
		currentAWSAccessKeyDefined = false
	}

	if envAWSSecretKey, found := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); found {
		currentAWSSecretKey = envAWSSecretKey
		currentAWSSecretKeyDefined = true
	} else {
		currentAWSSecretKeyDefined = false
	}

	if envAWSZoneId, found := os.LookupEnv("AWS_ZONE_ID"); found {
		currentAWSSecretKey = envAWSZoneId
		currentAWSZoneIdDefined = true
	} else {
		currentAWSZoneIdDefined = false
	}

	if envSubdomain, found := os.LookupEnv("SUBDOMAIN"); found {
		currentSubdomain = envSubdomain
		currentSubdomainDefined = true
	} else {
		currentSubdomainDefined = false
	}

	if envPowerDNSAPIHost, found := os.LookupEnv("POWER_DNS_API_HOST"); found {
		currentPowerDNSHost = envPowerDNSAPIHost
		currentPowerDNSHostDefined = true
	} else {
		currentPowerDNSHostDefined = false
	}

	if envPowerDNSAPIPort, found := os.LookupEnv("POWER_DNS_API_PORT"); found {
		currentPowerDNSPort = envPowerDNSAPIPort
		currentPowerDNSPortDefined = true
	} else {
		currentPowerDNSPortDefined = false
	}

	if envPowerDNSAPIKey, found := os.LookupEnv("POWER_DNS_API_KEY"); found {
		currentPowerDNSAPIKey = envPowerDNSAPIKey
		currentPowerDNSAPIKeyDefined = true
	} else {
		currentPowerDNSAPIKeyDefined = false
	}

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ZONE_ID")
	os.Unsetenv("SUBDOMAIN")
	os.Unsetenv("POWER_DNS_API_HOST")
	os.Unsetenv("POWER_DNS_API_PORT")
	os.Unsetenv("POWER_DNS_API_KEY")

}

func teardown() {

	if currentAWSAccessKeyDefined {
		os.Setenv("AWS_ACCESS_KEY_ID", currentAWSAccessKey)
	} else {
		os.Unsetenv("AWS_ACCESS_KEY_ID")
	}

	if currentAWSSecretKeyDefined {
		os.Setenv("AWS_SECRET_ACCESS_KEY", currentAWSSecretKey)
	} else {
		os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	}

	if currentAWSZoneIdDefined {
		os.Setenv("AWS_ZONE_ID", currentAWSZoneId)
	} else {
		os.Unsetenv("AWS_ZONE_ID")
	}

	if currentPowerDNSHostDefined {
		os.Setenv("POWER_DNS_API_HOST", currentPowerDNSHost)
	} else {
		os.Unsetenv("POWER_DNS_API_HOST")
	}

	if currentPowerDNSPortDefined {
		os.Setenv("POWER_DNS_API_PORT", currentPowerDNSPort)
	} else {
		os.Unsetenv("POWER_DNS_API_PORT")
	}

	if currentPowerDNSAPIKeyDefined {
		os.Setenv("POWER_DNS_API_KEY", currentPowerDNSAPIKey)
	} else {
		os.Unsetenv("POWER_DNS_API_KEY")
	}

}

func TestUpdaterWithInvalidAWSCredentails(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	updater := AWSUpdater{
		ZoneID:    "any",
		Subdomain: "any",
		IP:        "any",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidAWSCredentails should fail with invalid AWS_ACCESS_KEY_ID.")
	}

}

func TestUpdaterWithInvalidSecretKey(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	updater := AWSUpdater{
		ZoneID:    "any",
		Subdomain: "any",
		IP:        "any",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("CI_AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidSecretKey should fail with invalid AWS_ACCESS_KEY_ID.")
	}

}

func TestUpdaterWithInvalidZoneID(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	updater := AWSUpdater{
		ZoneID:    "any",
		Subdomain: "any",
		IP:        "any",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("CI_AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("CI_AWS_SECRET_ACCESS_KEY"))

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidAWSCredentails should fail with invalid AWS_ACCESS_KEY_ID.")
	}
}

func TestUpdaterWithInvalidPowerDNSAPIKey(t *testing.T) {

	setUp()
	defer teardown()

	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 seconds
	}

	ctx := context.TODO()
	api_port, _ := strconv.Atoi(os.Getenv("CI_POWER_DNS_API_PORT"))
	powerDNSClient, _ := powerdnsclient.NewClient(httpClient, os.Getenv("CI_POWER_DNS_API_HOST"), api_port, "invalidkey")

	updater := PowerDNSUpdater{
		PowerDNSClient: powerDNSClient,
		ZoneName:       "any",
		Subdomain:      "any",
		IP:             "any",
	}

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidPowerDNSAPIKey should fail with invalid PowerDNS API key.")
	}
}

func TestUpdaterWithValidPowerDNSAPIKey(t *testing.T) {

	setUp()
	defer teardown()

	httpClient := http.Client{
		Timeout: time.Second * 5, // Maximum of 5 seconds
	}

	ctx := context.TODO()
	api_port, _ := strconv.Atoi(os.Getenv("CI_POWER_DNS_API_PORT"))
	powerDNSClient, _ := powerdnsclient.NewClient(httpClient, os.Getenv("CI_POWER_DNS_API_HOST"), api_port, os.Getenv("CI_POWER_DNS_API_KEY"))

	updater := PowerDNSUpdater{
		PowerDNSClient: powerDNSClient,
		ZoneName:       "windmaker.net",
		Subdomain:      os.Getenv("CI_SUBDOMAIN"),
		IP:             "192.168.1.1",
	}

	err := updater.Update(ctx)

	if err != nil {
		t.Errorf("TestUpdaterWithValidPowerDNSAPIKey should not fail with valid PowerDNS API key. Error was \"%s\"", err.Error())
	}
}

func TestUpdaterWithValidData(t *testing.T) {

	ctx := context.TODO()

	updater := AWSUpdater{
		ZoneID:    os.Getenv("CI_ZONE_ID"),
		Subdomain: os.Getenv("CI_SUBDOMAIN"),
		IP:        "192.168.1.1",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("CI_AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("CI_AWS_SECRET_ACCESS_KEY"))
	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))

	err := updater.Update(ctx)

	if err != nil {
		t.Errorf("TestUpdaterWithValidData should not fail, error was \"%s\"", err.Error())
	}
}
