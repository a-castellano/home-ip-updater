package config

import (
	"cmp"
	"errors"
	"os"

	rabbitmqconfig "github.com/a-castellano/go-types/rabbitmq"
)

// Config struct contians required config variables
type Config struct {
	AWSZoneID      string // home-ip-monitor will send new IP values to be updated if associated ISP is the same than this value
	Subdomain      string // This will be the queue used to send IP changes
	UpdateQueue    string // This will be the queue used to send IP changes
	AWSRegion      string
	RabbitmqConfig *rabbitmqconfig.Config
}

// NewConfig checks if required env variables are present, returns config instance
func NewConfig() (*Config, error) {
	config := Config{}

	var envVariableFound bool
	// First check for AWS_ACCESS_KEY_ID env variable
	if _, envVariableFound = os.LookupEnv("AWS_ACCESS_KEY_ID"); !envVariableFound {
		return nil, errors.New("AWS_ACCESS_KEY_ID env variable must be set")
	}

	// Now check for AWS_SECRET_ACCESS_KEY
	if _, envVariableFound = os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !envVariableFound {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY env variable must be set")
	}

	// Above env variables are not required in config

	config.AWSRegion = cmp.Or(os.Getenv("AWS_REGION"), "us-west-2")

	if config.AWSZoneID, envVariableFound = os.LookupEnv("AWS_ZONE_ID"); !envVariableFound {
		return nil, errors.New("AWS_ZONE_ID env variable must be set")
	}

	if config.Subdomain, envVariableFound = os.LookupEnv("SUBDOMAIN"); !envVariableFound {
		return nil, errors.New("SUBDOMAIN env variable must be set")
	}

	var rabbitmqConfigErr error
	config.RabbitmqConfig, rabbitmqConfigErr = rabbitmqconfig.NewConfig()
	if rabbitmqConfigErr != nil {
		return nil, rabbitmqConfigErr
	}

	return &config, nil
}
