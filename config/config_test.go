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

	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ZONE_ID")
	os.Unsetenv("SUBDOMAIN")

	os.Unsetenv("RABBITMQ_HOST")
	os.Unsetenv("RABBITMQ_PORT")
	os.Unsetenv("RABBITMQ_DATABASE")
	os.Unsetenv("RABBITMQ_PASSWORD")
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

	_, err := NewConfig()

	if err == nil {
		t.Errorf("TestConfigWithRabbitmqInvalidPort should fail.")
	}

}

func TestValidConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")

	_, err := NewConfig()

	if err != nil {
		t.Errorf("TestConfigWithRabbitmqInvalidPort should not fail.")
	}

}
