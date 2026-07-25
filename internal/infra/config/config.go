// Package config provides configuration management for the home-ip-updater service.
// It handles environment variable validation and provides a centralized configuration
// structure for AWS and RabbitMQ settings.
package config

import (
	"cmp"
	"errors"
	"os"

	rabbitmqconfig "github.com/a-castellano/go-types/rabbitmq"
)

// Config contains all configuration variables required by the home-ip-updater service.
// It includes settings for AWS Route53, RabbitMQ, and domain configuration.
type Config struct {
	AWSZoneID      string                 // AWS Route53 hosted zone ID for DNS updates
	Subdomain      string                 // Subdomain to update with new IP addresses
	UpdateQueue    string                 // RabbitMQ queue name for receiving IP updates
	RabbitmqConfig *rabbitmqconfig.Config // RabbitMQ connection configuration
}

// NewConfig validates and loads all required environment variables into a Config struct.
// It performs validation for AWS credentials and RabbitMQ configuration.
// Returns an error if any required environment variables are missing or invalid.
//
// Required environment variables:
//   - AWS_ACCESS_KEY_ID: AWS access key for Route53 API access
//   - AWS_SECRET_ACCESS_KEY: AWS secret key for Route53 API access
//   - AWS_ZONE_ID: Route53 hosted zone ID
//   - SUBDOMAIN: Subdomain to update
//
// Optional environment variables:
//   - UPDATE_QUEUE_NAME: RabbitMQ queue name (defaults to "home-ip-monitor-updates")
//   - AWS_REGION: AWS region (defaults to "us-west-2")
func NewConfig() (*Config, error) {
	config := Config{}

	var envVariableFound bool

	// Validate AWS credentials - required for Route53 DNS updates
	if _, envVariableFound = os.LookupEnv("AWS_ACCESS_KEY_ID"); !envVariableFound {
		return nil, errors.New("AWS_ACCESS_KEY_ID env variable must be set")
	}

	if _, envVariableFound = os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !envVariableFound {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY env variable must be set")
	}

	// Set RabbitMQ queue name with default value
	config.UpdateQueue = cmp.Or(os.Getenv("UPDATE_QUEUE_NAME"), "home-ip-monitor-updates")

	// Set AWS region with default value
	AWSRegion := cmp.Or(os.Getenv("AWS_REGION"), "us-west-2")
	os.Setenv("AWS_REGION", AWSRegion)

	// Validate AWS Route53 configuration
	if config.AWSZoneID, envVariableFound = os.LookupEnv("AWS_ZONE_ID"); !envVariableFound {
		return nil, errors.New("AWS_ZONE_ID env variable must be set")
	}

	if config.Subdomain, envVariableFound = os.LookupEnv("SUBDOMAIN"); !envVariableFound {
		return nil, errors.New("SUBDOMAIN env variable must be set")
	}

	// Load RabbitMQ configuration from environment variables
	var rabbitmqConfigErr error
	config.RabbitmqConfig, rabbitmqConfigErr = rabbitmqconfig.NewConfig()
	if rabbitmqConfigErr != nil {
		return nil, rabbitmqConfigErr
	}

	return &config, nil
}
