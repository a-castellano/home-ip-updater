//go:build integration_tests || unit_tests

// Package config_test provides comprehensive testing for the config package.
// It includes unit tests for environment variable validation and configuration loading.
package config

import (
	"context"
	"os"
	"testing"
)

// envVariable holds the state of a single environment variable while the
// tests run: its name, the value it had before the suite touched it, and
// whether it was defined at all.
type envVariable struct {
	Value        string
	IsDefined    bool
	VariableName string
}

// envVariables lists every variable this package reads or writes, so setUp
// and teardown can save, clear and restore them generically. AWS_REGION
// belongs here even though it is optional: NewConfig writes it back with
// os.Setenv, so the suite would otherwise leak it into the environment.
//
// The env handling in this file (this map plus setUp and teardown) was
// rewritten by an AI agent (Claude); the test functions below are unchanged.
var envVariables = map[string]envVariable{
	//aws
	"aws_access_key": {VariableName: "AWS_ACCESS_KEY_ID"},
	"aws_secret_key": {VariableName: "AWS_SECRET_ACCESS_KEY"},
	"aws_zone_id":    {VariableName: "AWS_ZONE_ID"},
	//rabbitmq
	"rabbitmq_host":     {VariableName: "RABBITMQ_HOST"},
	"rabbitmq_port":     {VariableName: "RABBITMQ_PORT"},
	"rabbitmq_user":     {VariableName: "RABBITMQ_USER"},
	"rabbitmq_password": {VariableName: "RABBITMQ_PASSWORD"},
	//home-ip-updater
	"subdomain":         {VariableName: "SUBDOMAIN"},
	"update_queue_name": {VariableName: "UPDATE_QUEUE_NAME"},
}

// setUp saves the current environment variables and clears them for testing.
// This ensures that tests start with a clean environment and can properly
// test the validation logic without interference from existing environment variables.
func setUp() {

	for key, variable := range envVariables {

		if envValue, found := os.LookupEnv(variable.VariableName); found {
			variable.Value = envValue
			variable.IsDefined = true
		} else {
			variable.IsDefined = false
		}

		os.Unsetenv(variable.VariableName)

		envVariables[key] = variable
	}

}

// teardown restores the original environment variables after each test.
// This ensures that tests don't affect each other and the environment
// is returned to its original state.
func teardown() {

	for _, variable := range envVariables {
		if variable.IsDefined {
			os.Setenv(variable.VariableName, variable.Value)
		} else {
			os.Unsetenv(variable.VariableName)
		}
	}

}

// TestConfigWithoutEnvVariables verifies that NewConfig returns an error
// when no environment variables are set. This tests the basic validation
// that ensures required configuration is present.
func TestConfigWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err == nil {
		t.Errorf("TestConfigWithoutEnvVariables should fail.")
	} else {
		if err.Error() != "AWS_ACCESS_KEY_ID env variable must be set" {
			t.Errorf("TestConfigWithoutEnvVariables error should be \"AWS_ACCESS_KEY_ID env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

// TestConfigWithoutSecretKeyVariable verifies that NewConfig returns an error
// when AWS_SECRET_ACCESS_KEY is missing, even if other variables are set.
func TestConfigWithoutSecretKeyVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err == nil {
		t.Errorf("TestConfigWithoutEnvVariables should fail.")
	} else {
		if err.Error() != "AWS_SECRET_ACCESS_KEY env variable must be set" {
			t.Errorf("TestConfigWithoutEnvVariables error should be \"AWS_SECRET_ACCESS_KEY env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

// TestConfigWithoutZoneIdVariable verifies that NewConfig returns an error
// when AWS_ZONE_ID is missing, even if AWS credentials are set.
func TestConfigWithoutZoneIdVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err == nil {
		t.Errorf("TestConfigWithoutZoneIdVariable should fail.")
	} else {
		if err.Error() != "AWS_ZONE_ID env variable must be set" {
			t.Errorf("TestConfigWithoutZoneIdVariable error should be \"AWS_ZONE_ID env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

// TestConfigWithoutSubdomainVariable verifies that NewConfig returns an error
// when SUBDOMAIN is missing, even if AWS credentials and zone ID are set.
func TestConfigWithoutSubdomainVariable(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err == nil {
		t.Errorf("TestConfigWithoutSubdomainVariable should fail.")
	} else {
		if err.Error() != "SUBDOMAIN env variable must be set" {
			t.Errorf("TestConfigWithoutSubdomainVariable error should be \"SUBDOMAIN env variable must be set\" but it was \"%s\".", err.Error())
		}
	}

}

// TestConfigWithRabbitmqInvalidPort verifies that NewConfig returns an error
// when RabbitMQ port is set to an invalid value.
func TestConfigWithRabbitmqInvalidPort(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")
	os.Setenv("RABBITMQ_PORT", "invalidport")

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err == nil {
		t.Errorf("TestConfigWithRabbitmqInvalidPort should fail.")
	}

}

// TestValidConfig verifies that NewConfig succeeds when all required
// environment variables are properly set. This is the happy path test.
func TestValidConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")

	ctx := context.Background()
	_, err := NewConfig(ctx)

	if err != nil {
		t.Errorf("TestValidConfig should not fail.")
	}

}
