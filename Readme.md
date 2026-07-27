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
- **Structured Logging**: `slog` JSON logs written to stdout, collected by the journal when run under systemd
- **Observability**: Opt-in OpenTelemetry traces and metrics exported over OTLP
- **Configurable**: Environment-based configuration management
- **Production Ready**: Systemd service with security hardening
- **AWS Integration**: Automatic DNS record updates in Route53

## Prerequisites

- **Go 1.26+** for building and development
- **RabbitMQ Server** for message queuing
- **AWS Account** with Route53 access
- **Linux/Unix** system for production deployment

## Configuration

### Environment Variables

#### Required Variables

| Variable                | Description                            | Example                                    |
| ----------------------- | -------------------------------------- | ------------------------------------------ |
| `AWS_ACCESS_KEY_ID`     | AWS access key for Route53 API access  | `"AKIAIOSFODNN7EXAMPLE"`                   |
| `AWS_SECRET_ACCESS_KEY` | AWS secret key for Route53 API access  | `"wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"` |
| `AWS_ZONE_ID`           | Route53 hosted zone ID for DNS updates | `"Z1234567890ABC"`                         |
| `SUBDOMAIN`             | Subdomain to update                    | `"home.example.com"`                       |

The Route53 client always talks to the `us-east-1` endpoint, since Route53 is a global service. No AWS region variable is read.

#### Optional Variables

| Variable            | Description                        | Default                     |
| ------------------- | ---------------------------------- | --------------------------- |
| `UPDATE_QUEUE_NAME` | RabbitMQ queue name for IP updates | `"home-ip-monitor-updates"` |

#### Application and Logging

Logging is handled through [go-types `slog`](https://git.windmaker.net/a-castellano/go-types/-/tree/master/slog). `APP_NAME` is required by that type; the rest fall back to sane defaults.

| Variable          | Description                                   | Default      |
| ----------------- | --------------------------------------------- | ------------ |
| `APP_NAME`        | Application name attached to every log entry  | _(required)_ |
| `SLOG_LEVEL`      | Log level: `Debug`, `Info`, `Warn` or `Error` | `Info`       |
| `SLOG_FORMAT`     | Log format: `JSON` or `plain`                 | `JSON`       |
| `SLOG_ADD_SOURCE` | Whether to add `file:line` to log entries     | `true`       |

#### Telemetry

OpenTelemetry is opt-in through [go-types `opentelemetry`](https://git.windmaker.net/a-castellano/go-types/-/tree/master/opentelemetry). `APP_NAME` doubles as the telemetry `service.name`, so `OTEL_SERVICE_NAME` and `OTEL_RESOURCE_ATTRIBUTES` must **not** be set — the config rejects them.

| Variable                      | Description                                                                                    | Default            |
| ----------------------------- | ---------------------------------------------------------------------------------------------- | ------------------ |
| `ENABLE_TELEMETRY`            | Enables traces and metrics when set to `"true"` (opt-in)                                       | `false`            |
| `OTEL_EXPORTER_OTLP_ENDPOINT` | OTLP collector endpoint (`http://` or `https://`); when unset, traces/metrics export to stdout | _(unset → stdout)_ |

#### RabbitMQ Configuration

The following RabbitMQ environment variables are required (see [go-types documentation](https://git.windmaker.net/a-castellano/go-types/-/tree/master/rabbitmq?ref_type=heads)):

| Variable            | Description              | Default       |
| ------------------- | ------------------------ | ------------- |
| `RABBITMQ_HOST`     | RabbitMQ server hostname | `"localhost"` |
| `RABBITMQ_PORT`     | RabbitMQ server port     | `5672`        |
| `RABBITMQ_USER`     | RabbitMQ username        | `"guest"`     |
| `RABBITMQ_PASSWORD` | RabbitMQ password        | `"guest"`     |

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
podman compose -f development/docker-compose.yml up -d

# Access services:
# - RabbitMQ Management: http://localhost:15672
# - RabbitMQ AMQP: localhost:5672
# - OpenTelemetry collector zPages: http://localhost:55679/debug/tracez
```

Go commands run inside the `golang` container, which uses the same image as CI and production:

```bash
podman compose -f development/docker-compose.yml exec golang make test
```

`development/env_variables` holds the environment used when running the service by hand against that environment.

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

The unit writes its logs to stdout and stderr, both collected by the journal, so `journalctl` is the only place to look for them.

### Environment Configuration for Systemd

The package installs a sample file at `/etc/default/windmaker-home-ip-updater-example`. Copy it to `/etc/default/windmaker-home-ip-updater` (the path read by the systemd unit) and edit it to configure the service:

```bash
sudo cp /etc/default/windmaker-home-ip-updater-example /etc/default/windmaker-home-ip-updater
sudo vim /etc/default/windmaker-home-ip-updater
```

```bash
# Application and logging
APP_NAME="home-ip-updater"
SLOG_LEVEL="Info"
SLOG_FORMAT="JSON"

# Telemetry (opt-in)
ENABLE_TELEMETRY="false"
# OTEL_EXPORTER_OTLP_ENDPOINT="http://otelcollector:4317"

# Queue configuration
UPDATE_QUEUE_NAME="home-ip-monitor-updates"

# AWS config
AWS_ACCESS_KEY_ID="your_aws_access_key"
AWS_SECRET_ACCESS_KEY="your_aws_secret_key"
AWS_ZONE_ID="your_route53_zone_id"
SUBDOMAIN="home.example.com"

# RabbitMQ config
RABBITMQ_HOST="localhost"
RABBITMQ_PORT=5672
RABBITMQ_USER="guest"
RABBITMQ_PASSWORD="guest"
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
