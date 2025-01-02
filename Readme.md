# Home IP Updater

[![pipeline status](https://git.windmaker.net/a-castellano/home-ip-updater/badges/master/pipeline.svg)](https://git.windmaker.net/a-castellano/home-ip-updater/pipelines)[![coverage report](https://git.windmaker.net/a-castellano/home-ip-updater/badges/master/coverage.svg)](https://a-castellano.gitpages.windmaker.net/home-ip-updater/coverage.html)[![Quality Gate Status](https://sonarqube.windmaker.net/api/project_badges/measure?project=a-castellano_home-ip-updater_533a7009-26fb-43b9-b6f3-eb5326c083b6&metric=alert_status&token=sqb_df6b40224599cede55c63c9203eb5fcdb0a4bc9e)](https://sonarqube.windmaker.net/dashboard?id=a-castellano_home-ip-updater_533a7009-26fb-43b9-b6f3-eb5326c083b6)

This program is subscribed to [home-ip-monitor](https://git.windmaker.net/a-castellano/home-ip-monitor) update queue, it will update required DNS record with readed IP's from queue.

# What this utility does?

Reads IP's from configured queue and updates required DNS record in AWS Route53 service and my local DNS server.

# Required variables

## Queue names

**UPDATE_QUEUE_NAME**: Queue name where new IP's will be sent.

## AWS Config

DNS service is hosted in AWS Route53 service, AWS auth and Route53 keys are required:

**AWS_ACCESS_KEY_ID**: AWS account key
**AWS_SECRET_ACCESS_KEY**: AWS account secret key
**AWS_REGION**: AWS region, default value is "us-west-2"
**AWS_ZONE_ID**: Route53 zone ID

In this case, an user with limited permissions is created with the following policy:

```
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

## PowerDNS config

This program uses PowerDSN API for updating records, the following env vars must be set:

**POWER_DNS_API_HOST**: PowerDNS API IP
**POWER_DNS_API_PORT**: PowerDNS API Port
**POWER_DNS_API_KEY**: API key used to access PowerDNS API
**POWER_DNS_ZONE_NAME**: DNS Zone name to update

## Domain config

Required domain A record will be updated, domain value is given by the following env variable:
**SUBDOMAIN**: Subdomain record to update

## RabbitMQ Config

RabbitMQ required config can be found in its [go types](https://git.windmaker.net/a-castellano/go-types/-/tree/master/rabbitmq?ref_type=heads) Readme.
