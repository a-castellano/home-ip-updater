# Home IP Updater

[![pipeline status](https://git.windmaker.net/a-castellano/home-ip-updater/badges/master/pipeline.svg)](https://git.windmaker.net/a-castellano/home-ip-updater/pipelines)[![coverage report](https://git.windmaker.net/a-castellano/home-ip-updater/badges/master/coverage.svg)](https://a-castellano.gitpages.windmaker.net/home-ip-updater/coverage.html)[![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=a-castellano_home-ip-updater_533a7009-26fb-43b9-b6f3-eb5326c083b6&metric=alert_status&token=sqb_df6b40224599cede55c63c9203eb5fcdb0a4bc9e)](https://sonarqube.windmaker.net/dashboard?id=a-castellano_home-ip-updater_533a7009-26fb-43b9-b6f3-eb5326c083b6)

Go microservice that monitors IP changes and updates DNS records in AWS Route53. This service is part of the home-ip-monitor ecosystem and subscribes to RabbitMQ queues to receive IP change notifications.

## What This Program Does

The Home IP Updater is designed to:

- **Subscribe** to a RabbitMQ queue for IP change notifications
- **Process** incoming IP change messages
- **Update** DNS records in AWS Route53 when IP changes are detected
- **Provide** reliable, production-ready DNS management for dynamic IP addresses

## Architecture

```
┌─────────────────┐    ┌──────────────────┐    ┌─────────────────┐
│  home-ip-monitor│───▶│  RabbitMQ Queue  │───▶│ home-ip-updater │
│  (IP Detector)  │    │                  │    │                 │
└─────────────────┘    └──────────────────┘    └─────────────────┘
                                                          │
                                                          ▼
                                               ┌─────────────────┐
                                               │   AWS Route53   │
                                               │   (DNS)         │
                                               └─────────────────┘
```

## Features

- **Reliable Message Processing**: Handles RabbitMQ messages with error recovery
- **Graceful Shutdown**: Proper signal handling for clean service termination
- **System Integration**: Syslog logging for production environments
- **Configurable**: Environment-based configuration management
- **Production Ready**: Systemd service with security hardening
- **AWS Integration**: Automatic DNS record updates in Route53

## Prerequisites

- **Go 1.24+** for building and development
- **RabbitMQ Server** for message queuing
- **AWS Account** with Route53 access
- **Linux/Unix** system for production deployment

## Configuration

### Environment Variables

#### Required Variables

| Variable              | Description                            | Example                                    |
| --------------------- | -------------------------------------- | ------------------------------------------ |
| AWS_ACCESS_KEY_ID     | AWS access key for Route53 API access  | "AKIAIOSFODNN7EXAMPLE"                     |
| AWS_SECRET_ACCESS_KEY | AWS secret key for Route53 API access  | "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY" |
| AWS_ZONE_ID           | Route53 hosted zone ID for DNS updates | "Z1234567890ABC"                           |
| SUBDOMAIN             | Subdomain to update                    | "home.example.com"                         |

#### Optional Variables

| Variable          | Description                        | Default                   |
| ----------------- | ---------------------------------- | ------------------------- |
| UPDATE_QUEUE_NAME | RabbitMQ queue name for IP updates | "home-ip-monitor-updates" |
| AWS_REGION        | AWS region for Route53 operations  | "us-west-2"               |

#### RabbitMQ Configuration

The following RabbitMQ environment variables are required (see [go-types documentation](https://git.windmaker.net/a-castellano/go-types/-/tree/master/rabbitmq?ref_type=heads)):

| Variable          | Description              | Default     |
| ----------------- | ------------------------ | ----------- |
| RABBITMQ_HOST     | RabbitMQ server hostname | "localhost" |
| RABBITMQ_PORT     | RabbitMQ server port     | "5672"      |
| RABBITMQ_USER     | RabbitMQ username        | "guest"     |
| RABBITMQ_PASSWORD | RabbitMQ password        | "guest"     |

### AWS IAM Policy

Create an IAM user with the following policy for Route53 access:

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "arn:aws:route53:::hostedzone/AWS_ZONE_ID"
    }
  ]
}
```

## Development

### Running Tests

```bash
# Run unit tests
make test

# Run integration tests
make test_integration

# Generate coverage report
make coverage

# Check code quality
make lint
```

### Development Environment

The project includes a Docker Compose setup for development:

```bash
# Start development environment
cd development
docker-compose up -d

# Access services:
# - RabbitMQ Management: http://localhost:15672
# - RabbitMQ AMQP: localhost:5672
```

### Building

```bash
# Build the service
make build

# Clean previous builds
make clean
```

## Deployment

### Systemd Service Installation

The application includes a systemd service file for production deployment:

```bash
# Install the service
sudo systemctl enable windmaker-home-ip-updater.service
sudo systemctl start windmaker-home-ip-updater.service

# Check service status
sudo systemctl status windmaker-home-ip-updater.service

# View logs
sudo journalctl -u windmaker-home-ip-updater.service -f
```

### Environment Configuration for Systemd

Create an environment file for the systemd service:

```bash
# Create environment file
sudo mkdir -p /etc/windmaker-home-ip-updater
sudo nano /etc/windmaker-home-ip-updater/environment
```

Add the following environment variables:

```bash
# AWS Configuration
AWS_ACCESS_KEY_ID=your_aws_access_key
AWS_SECRET_ACCESS_KEY=your_aws_secret_key
AWS_ZONE_ID=your_route53_zone_id
AWS_REGION=us-west-2

# Domain Configuration
SUBDOMAIN=home.example.com

# RabbitMQ Configuration
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest

# Queue Configuration
UPDATE_QUEUE_NAME=home-ip-monitor-updates
```

## License

This project is licensed under the GPL v3 License - see the LICENSE file for details.

## Authors

- **Álvaro Castellano Vela** - _Initial work_ - a-castellano

## Related Projects

- [home-ip-monitor](https://git.windmaker.net/a-castellano/home-ip-monitor) - The IP monitoring service that feeds this updater
- [home-ip-notifier](https://github.com/a-castellano/home-ip-notifier) - Email notification service for IP changes
- [go-services](https://git.windmaker.net/a-castellano/go-services) - Shared Go service utilities
- [go-types](https://git.windmaker.net/a-castellano/go-types) - Shared Go type definitions
