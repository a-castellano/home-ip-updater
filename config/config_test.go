//go:build integration_tests || unit_tests

package config

import (
	"os"
	"testing"
)

var currentAWSAccessKey string
var currentAWSAccessKeyDefined bool

var currentAWSSecretKey string
var currentAWSSecretKeyDefined bool

var currentAWSZoneId string
var currentAWSZoneIdDefined bool

var currentSubdomain string
var currentSubdomainDefined bool

var currentRabbitmqHost string
var currentRabbitmqHostDefined bool

var currentRabbitmqPort string
var currentRabbitmqPortDefined bool

var currentRabbitmqUser string
var currentRabbitmqUserDefined bool

var currentRabbitmqPassword string
var currentRabbitmqPasswordDefined bool

var currentPowerDNSHost string
var currentPowerDNSHostDefined bool

var currentPowerDNSPort string
var currentPowerDNSPortDefined bool

var currentPowerDNSAPIKey string
var currentPowerDNSAPIKeyDefined bool

var currentPowerDNSZoneName string
var currentPowerDNSZoneNameDefined bool

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

	if envPowerDNSZoneName, found := os.LookupEnv("POWER_DNS_ZONE_NAME"); found {
		currentPowerDNSZoneName = envPowerDNSZoneName
		currentPowerDNSZoneNameDefined = true
	} else {
		currentPowerDNSZoneNameDefined = false
	}

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ZONE_ID")
	os.Unsetenv("SUBDOMAIN")

	os.Unsetenv("RABBITMQ_HOST")
	os.Unsetenv("RABBITMQ_PORT")
	os.Unsetenv("RABBITMQ_DATABASE")
	os.Unsetenv("RABBITMQ_PASSWORD")

	os.Unsetenv("POWER_DNS_API_HOST")
	os.Unsetenv("POWER_DNS_API_PORT")
	os.Unsetenv("POWER_DNS_API_KEY")
	os.Unsetenv("POWER_DNS_ZONE_NAME")
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

	if currentRabbitmqHostDefined {
		os.Setenv("RABBITMQ_HOST", currentRabbitmqHost)
	} else {
		os.Unsetenv("RABBITMQ_HOST")
	}

	if currentRabbitmqPortDefined {
		os.Setenv("RABBITMQ_PORT", currentRabbitmqPort)
	} else {
		os.Unsetenv("RABBITMQ_PORT")
	}

	if currentRabbitmqUserDefined {
		os.Setenv("RABBITMQ_USER", currentRabbitmqUser)
	} else {
		os.Unsetenv("RABBITMQ_USER")
	}

	if currentRabbitmqPasswordDefined {
		os.Setenv("RABBITMQ_PASSWORD", currentRabbitmqPassword)
	} else {
		os.Unsetenv("RABBITMQ_PASSWORD")
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

	if currentPowerDNSZoneNameDefined {
		os.Setenv("POWER_DNS_ZONE_NAME", currentPowerDNSZoneName)
	} else {
		os.Unsetenv("POWER_DNS_ZONE_NAME")
	}
}

func TestConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutEnvVariables should fail.")
	} else {
		if err.Error() != "AWS_ACCESS_KEY_ID env variable must be set" {
			t.Errorf("TestConfigWithoutEnvVariables error should be \"AWS_ACCESS_KEY_ID env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

func TestConfigWithoutSecretKeyVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutEnvVariables should fail.")
	} else {
		if err.Error() != "AWS_SECRET_ACCESS_KEY env variable must be set" {
			t.Errorf("TestConfigWithoutEnvVariables error should be \"AWS_SECRET_ACCESS_KEY env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

func TestConfigWithoudZoneIdariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_PORT", "8080")
	os.Setenv("POWER_DNS_API_KEY", "key")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoudZoneIdariable gshould fail.")
	} else {
		if err.Error() != "AWS_ZONE_ID env variable must be set" {
			t.Errorf("TestConfigWithoutSecretKeyVariable error should be \"AWS_ZONE_ID env variable must be set env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

func TestConfigWithoudSubdomainariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_PORT", "8080")
	os.Setenv("POWER_DNS_API_KEY", "key")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoudZoneIdashould fail.")
	} else {
		if err.Error() != "SUBDOMAIN env variable must be set" {
			t.Errorf("TestConfigWithoutSecretKeyVariable error should be \"SUBDOMAIN env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

func TestConfigWithRabbitmqInvalidPort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("RABBITMQ_PORT", "invalidport")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithRabbitmqInvalidPort should fail.")
	}

}

func TestConfigWithoutPowerDNSAPIHost(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("POWER_DNS_API_PORT", "8080")
	os.Setenv("POWER_DNS_API_KEY", "key")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutPowerDNSAPIHost should fail.")
	}

}

func TestConfigWithoutPowerDNSAPIPort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_KEY", "key")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutPowerDNSAPIPort should fail.")
	}

}

func TestConfigWithoutPowerDNSAPIKey(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_PORT", "8080")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutPowerDNSAPIKey should fail.")
	}

}

func TestConfigWithoutPowerDNSZoneName(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_PORT", "8080")

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithoutPowerDNSAPIKey should fail.")
	}

}

func TestValidConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("POWER_DNS_API_HOST", "host")
	os.Setenv("POWER_DNS_API_PORT", "8080")
	os.Setenv("POWER_DNS_API_KEY", "key")
	os.Setenv("POWER_DNS_ZONE_NAME", "test.net")

	_, err := NewConfig()

	if err != nil {
		t.Errorf("TestConfigWithRabbitmqInvalidPort should not fail.")
	}

}
