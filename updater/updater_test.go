//go:build integration_tests || unit_tests

package updater

import (
	"context"
	"fmt"
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
		t.Errorf("TestUpdaterWithValidData should nor fail, error was \"%s\"", err.Error())
	}
}
