//go:build integration_tests

package route53

import (
	"context"
	"os"
	"strings"
	"testing"

	appconfig "github.com/a-castellano/home-ip-updater/internal/infra/config"
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

func TestUpdaterWithoutEnvVariables(t *testing.T) {

	setUp()
	defer teardown()

	config := appconfig.Config{
		AWSZoneID: "invalid",
		AWSRegion: "invalid",
		Subdomain: "invalid",
	}

	ctx := context.Background()

	_, err := NewRoute53Updater(ctx, &config)

	if err == nil {
		t.Fatalf("TestUpdaterWithoutEnvVariables should fail.")
	}

	expectedError := "no EC2 IMDS role found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Fatalf("TestUpdaterWithoutEnvVariables error should contain \"%s\", it was \"%s\"", expectedError, err.Error())
	}
}
