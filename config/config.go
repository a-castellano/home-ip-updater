package config

import (
	"cmp"
	"errors"
	"os"
	"strconv"

	rabbitmqconfig "github.com/a-castellano/go-types/rabbitmq"
)

// Config struct contians required config variables
type Config struct {
	AWSZoneID        string // home-ip-monitor will send new IP values to be updated if associated ISP is the same than this value
	Subdomain        string // Subdomain to update
	PowerDNSAPIKey   string // PowerDNS API key
	PowerDNSHost     string // PowerDNS API Host
	PowerDNSPort     int    // PowerDNS API Port
	PowerDNSZoneName string // PowerDNS API Port
	UpdateQueue      string // This will be the queue used to send IP changes
	RabbitmqConfig   *rabbitmqconfig.Config
}

// NewConfig checks if required env variables are present, returns config instance
func NewConfig() (*Config, error) {
	config := Config{}

	var envVariableFound bool
	var PowerDNSPortString string
	// First check for AWS_ACCESS_KEY_ID env variable
	if _, envVariableFound = os.LookupEnv("AWS_ACCESS_KEY_ID"); !envVariableFound {
		return nil, errors.New("AWS_ACCESS_KEY_ID env variable must be set")
	}

	// Now check for AWS_SECRET_ACCESS_KEY
	if _, envVariableFound = os.LookupEnv("AWS_SECRET_ACCESS_KEY"); !envVariableFound {
		return nil, errors.New("AWS_SECRET_ACCESS_KEY env variable must be set")
	}

	// All PowerDNS variables are required
	if config.PowerDNSHost, envVariableFound = os.LookupEnv("POWER_DNS_API_HOST"); !envVariableFound {
		return nil, errors.New("POWER_DNS_API_HOST env variable must be set")
	}
	if PowerDNSPortString, envVariableFound = os.LookupEnv("POWER_DNS_API_PORT"); !envVariableFound {
		return nil, errors.New("POWER_DNS_API_PORT env variable must be set")
	}
	config.PowerDNSPort, _ = strconv.Atoi(PowerDNSPortString)
	if config.PowerDNSAPIKey, envVariableFound = os.LookupEnv("POWER_DNS_API_KEY"); !envVariableFound {
		return nil, errors.New("POWER_DNS_API_KEY env variable must be set")
	}
	if config.PowerDNSZoneName, envVariableFound = os.LookupEnv("POWER_DNS_ZONE_NAME"); !envVariableFound {
		return nil, errors.New("POWER_DNS_ZONE_NAME env variable must be set")
	}

	// Above env variables are not required in config

	// Retrieve UpdateQueue name, default is home-ip-monitor-updates
	config.UpdateQueue = cmp.Or(os.Getenv("UPDATE_QUEUE_NAME"), "home-ip-monitor-updates")

	// Sets AWS region
	AWSRegion := cmp.Or(os.Getenv("AWS_REGION"), "us-west-2")
	os.Setenv("AWS_REGION", AWSRegion)

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
