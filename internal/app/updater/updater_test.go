//go:build integration_tests || unit_tests

// Package updater_test provides comprehensive testing for the updater package.
// It includes unit tests for AWS Route53 DNS record updates and error handling.
package updater

import (
	"context"
	//	"fmt"
	"os"
	"testing"
)

// Global variables to store original environment variable values
// These are used to restore the environment state after tests
var currentAWSAccessKey string
var currentAWSAccessKeyDefined bool

var currentAWSSecretKey string
var currentAWSSecretKeyDefined bool

var currentAWSZoneId string
var currentAWSZoneIdDefined bool

var currentSubdomain string
var currentSubdomainDefined bool

// setUp saves the current environment variables and clears them for testing.
// This ensures that tests start with a clean environment and can properly
// test AWS credential validation without interference from existing environment variables.
func setUp() {

	// Save AWS access key if it exists
	if envAWSAccessKey, found := os.LookupEnv("AWS_ACCESS_KEY_ID"); found {
		currentAWSAccessKey = envAWSAccessKey
		currentAWSAccessKeyDefined = true
	} else {
		currentAWSAccessKeyDefined = false
	}

	// Save AWS secret key if it exists
	if envAWSSecretKey, found := os.LookupEnv("AWS_SECRET_ACCESS_KEY"); found {
		currentAWSSecretKey = envAWSSecretKey
		currentAWSSecretKeyDefined = true
	} else {
		currentAWSSecretKeyDefined = false
	}

	// Save AWS zone ID if it exists
	if envAWSZoneId, found := os.LookupEnv("AWS_ZONE_ID"); found {
		currentAWSZoneId = envAWSZoneId
		currentAWSZoneIdDefined = true
	} else {
		currentAWSZoneIdDefined = false
	}

	// Save subdomain if it exists
	if envSubdomain, found := os.LookupEnv("SUBDOMAIN"); found {
		currentSubdomain = envSubdomain
		currentSubdomainDefined = true
	} else {
		currentSubdomainDefined = false
	}

	// Clear all environment variables to ensure clean test state
	os.Unsetenv("AWS_ACCESS_KEY_ID")
	os.Unsetenv("AWS_SECRET_ACCESS_KEY")
	os.Unsetenv("AWS_ZONE_ID")
	os.Unsetenv("SUBDOMAIN")

}

// teardown restores the original environment variables after each test.
// This ensures that tests don't affect each other and the environment
// is returned to its original state.
func teardown() {

	// Restore AWS credentials if they existed before
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

}

// TestUpdaterWithInvalidAWSCredentials verifies that the AWS updater returns an error
// when invalid AWS credentials are provided. This tests the AWS SDK's credential
// validation and error handling.
func TestUpdaterWithInvalidAWSCredentials(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	// Create updater with invalid credentials
	updater := AWSUpdater{
		ZoneID:    "any",
		Subdomain: "any",
		IP:        "any",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "test")

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidAWSCredentials should fail with invalid AWS_ACCESS_KEY_ID.")
	}

}

// TestUpdaterWithInvalidSecretKey verifies that the AWS updater returns an error
// when a valid access key but invalid secret key is provided. This tests
// the AWS SDK's credential validation for secret keys.
func TestUpdaterWithInvalidSecretKey(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	// Create updater with valid access key but invalid secret key
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

// TestUpdaterWithInvalidZoneID verifies that the AWS updater returns an error
// when an invalid hosted zone ID is provided. This tests the Route53 API's
// validation of hosted zone IDs and error handling.
func TestUpdaterWithInvalidZoneID(t *testing.T) {

	setUp()
	defer teardown()

	ctx := context.TODO()

	// Create updater with invalid zone ID
	updater := AWSUpdater{
		ZoneID:    "any",
		Subdomain: "any",
		IP:        "any",
	}
	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("CI_AWS_ACCESS_KEY_ID"))
	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("CI_AWS_SECRET_ACCESS_KEY"))

	err := updater.Update(ctx)

	if err == nil {
		t.Errorf("TestUpdaterWithInvalidAWSCredentials should fail with invalid AWS_ACCESS_KEY_ID.")
	}
}

//// TestUpdaterWithValidData verifies that the AWS updater succeeds when all
//// parameters are valid. This is the happy path test that ensures the DNS
//// update functionality works correctly with proper credentials and configuration.
//// This test requires valid AWS credentials and a real hosted zone ID to pass.
//func TestUpdaterWithValidData(t *testing.T) {
//
//	ctx := context.TODO()
//
//	// Create updater with valid configuration
//	updater := AWSUpdater{
//		ZoneID:    os.Getenv("CI_ZONE_ID"),
//		Subdomain: os.Getenv("CI_SUBDOMAIN"),
//		IP:        "192.168.1.1",
//	}
//	os.Setenv("AWS_ACCESS_KEY_ID", os.Getenv("CI_AWS_ACCESS_KEY_ID"))
//	os.Setenv("AWS_SECRET_ACCESS_KEY", os.Getenv("CI_AWS_SECRET_ACCESS_KEY"))
//	fmt.Println(os.Getenv("AWS_SECRET_ACCESS_KEY"))
//
//	err := updater.Update(ctx)
//
//	if err != nil {
//		t.Errorf("TestUpdaterWithValidData should not fail, error was \"%s\"", err.Error())
//	}
//}
