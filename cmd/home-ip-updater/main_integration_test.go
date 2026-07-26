//go:build integration_tests

package main

import (
	"context"
	"encoding/json"
	rabbitmq "github.com/a-castellano/go-services/infra/rabbitmq"
	messagebroker "github.com/a-castellano/go-services/services/messagebroker"
	envelope "github.com/a-castellano/go-types/types/envelope"
	rabbitmqconfig "github.com/a-castellano/go-types/types/rabbitmq"
	"net/http"
	"os"
	"testing"
	"time"
)

type envVariable struct {
	Value        string
	IsDefined    bool
	VariableName string
}

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
	"app_name":          {VariableName: "APP_NAME"},
}

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

func teardown() {

	for _, variable := range envVariables {
		if variable.IsDefined {
			os.Setenv(variable.VariableName, variable.Value)
		} else {
			os.Unsetenv(variable.VariableName)
		}
	}
}

func TestInvalidOtelConfig(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")
	os.Setenv("SUBDOMAIN", "test.windmaker.net")

	ctx := context.Background()
	err := run(ctx)

	if err == nil {
		t.Fatalf("TestInvalidOtelConfig should fail, otel config is invalid")
	}

}

func TestNoSubdomain(t *testing.T) {

	setUp()
	defer teardown()

	os.Setenv("APP_NAME", "home-ip-updater")

	os.Setenv("AWS_ACCESS_KEY_ID", "test")
	os.Setenv("AWS_SECRET_ACCESS_KEY", "secret")
	os.Setenv("AWS_ZONE_ID", "123")

	ctx := context.Background()
	err := run(ctx)

	if err == nil {
		t.Fatalf("TestNoSubdomain should fail as subdomain is not set")
	}

	expectedErr := "SUBDOMAIN env variable must be set"
	if err.Error() != expectedErr {
		t.Fatalf("TestNoSubdomain error should be \"%s\" but it was \"%s\".", expectedErr, err.Error())
	}

}

// Secrets mirrors the contents of development/secrets.json, the git-crypt
// encrypted file that also feeds the dns package integration tests. run needs
// real AWS credentials because NewRoute53Updater validates them against
// Route53 before the service reaches its message loop.
type Secrets struct {
	AWSAccessKeyID     string `json:"AWS_ACCESS_KEY_ID"`
	AWSSecretAccessKey string `json:"AWS_SECRET_ACCESS_KEY"`
	AWSZoneID          string `json:"AWS_ZONE_ID"`
	Subdomain          string `json:"SUBDOMAIN"`
}

const rabbitmqAPI = "http://rabbitmq:15672"
const integrationQueueName = "home-ip-updater-integration-tests"

// queueIsDrained reports whether queueName has a consumer attached and no
// pending messages, which is how this suite observes that the service picked
// the test message up. Anything other than a 200 means the management API has
// no stats for the queue yet, so it counts as "not drained" instead of an
// error and the caller keeps polling.
// This helper was written by an AI agent (Claude).
func queueIsDrained(queueName string) (bool, error) {
	request, requestErr := http.NewRequest(http.MethodGet, rabbitmqAPI+"/api/queues/%2F/"+queueName, nil)
	if requestErr != nil {
		return false, requestErr
	}
	request.SetBasicAuth("guest", "guest")

	response, getErr := http.DefaultClient.Do(request)
	if getErr != nil {
		return false, getErr
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return false, nil
	}

	var queue struct {
		Messages  int `json:"messages"`
		Consumers int `json:"consumers"`
	}
	if decodeErr := json.NewDecoder(response.Body).Decode(&queue); decodeErr != nil {
		return false, decodeErr
	}

	return queue.Consumers > 0 && queue.Messages == 0, nil
}

// TestUpdateIP publishes an IP update on the queue the service consumes and
// checks that run picks it up, stays alive, and shuts down without error once
// its context is cancelled. It asserts nothing about Route53: the DNS side is
// already covered by the dns package integration tests, so the observable
// outcome here is the queue being drained by the service's own consumer.
//
// This test was written by an AI agent (Claude), adapting TestRetrieveMessage
// from home-ip-notifier: the MailHog polling of the original is replaced by
// RabbitMQ management API polling, and the SMTP setup by the AWS credentials
// run needs to build its Route53 updater.
func TestUpdateIP(t *testing.T) {

	setUp()
	defer teardown()

	var secrets Secrets
	secretData, readErr := os.ReadFile("../../development/secrets.json")
	if readErr != nil {
		t.Fatalf("TestUpdateIP should not fail reading secret file, error was \"%s\"", readErr.Error())
	}

	if jsonErr := json.Unmarshal(secretData, &secrets); jsonErr != nil {
		t.Fatalf("TestUpdateIP should not fail reading json file content, error was \"%s\"", jsonErr.Error())
	}

	os.Setenv("APP_NAME", "home-ip-updater")

	os.Setenv("AWS_ACCESS_KEY_ID", secrets.AWSAccessKeyID)
	os.Setenv("AWS_SECRET_ACCESS_KEY", secrets.AWSSecretAccessKey)
	os.Setenv("AWS_ZONE_ID", secrets.AWSZoneID)
	os.Setenv("SUBDOMAIN", secrets.Subdomain)

	// Set environment variables for RabbitMQ with valid credentials
	os.Setenv("RABBITMQ_HOST", "rabbitmq")
	os.Setenv("RABBITMQ_PORT", "5672")
	os.Setenv("RABBITMQ_USER", "guest")
	os.Setenv("RABBITMQ_PASSWORD", "guest")

	os.Setenv("UPDATE_QUEUE_NAME", integrationQueueName)

	rabbitmqConfig, rabbitmqConfigErr := rabbitmqconfig.NewConfig()
	if rabbitmqConfigErr != nil {
		t.Fatalf("TestUpdateIP could not load rabbitmq config, error was \"%s\"", rabbitmqConfigErr.Error())
	}

	rabbitmqClient := rabbitmq.NewRabbitmqClient(rabbitmqConfig)
	messageBroker := messagebroker.MessageBroker{Client: rabbitmqClient}

	body := []byte("123.123.123.123")
	carrier := map[string]string{"traceparent": "00-0af7651916cd43dd8448eb211c80319c-b7ad6b7169203331-01"}

	data, marshalErr := (&envelope.Envelope{Carrier: carrier, Body: body}).Marshal()
	if marshalErr != nil {
		t.Fatalf("TestUpdateIP could not marshal the test envelope, error was \"%s\"", marshalErr.Error())
	}

	sendError := messageBroker.SendMessage(context.Background(), integrationQueueName, data)

	if sendError != nil {
		t.Fatalf("TestUpdateIP should not fail when test message is sent, error was \"%s\"", sendError.Error())
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	runResult := make(chan error, 1)
	go func() {
		runResult <- run(ctx)
	}()

	// Poll the management API until the queue is drained; run exiting early or
	// the deadline expiring are both failures. The deadline is generous
	// because RabbitMQ refreshes these statistics every few seconds.
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()
	deadline := time.After(45 * time.Second)

	consumed := false
	for !consumed {
		select {
		case runErr := <-runResult:
			t.Fatalf("TestUpdateIP run returned before the message was consumed, error was \"%v\"", runErr)
		case <-deadline:
			t.Fatalf("TestUpdateIP timed out waiting for the service to consume the test message")
		case <-ticker.C:
			drained, checkErr := queueIsDrained(integrationQueueName)
			if checkErr != nil {
				t.Fatalf("TestUpdateIP could not query the RabbitMQ management API, error was \"%s\"", checkErr.Error())
			}
			consumed = drained
		}
	}

	cancel()

	select {
	case runErr := <-runResult:
		if runErr != nil {
			t.Fatalf("TestUpdateIP run should finish without error after cancellation, error was \"%s\"", runErr.Error())
		}
	case <-time.After(5 * time.Second):
		t.Fatalf("TestUpdateIP run did not return after context cancellation")
	}

}
