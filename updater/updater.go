// Package updater provides DNS record update functionality for the home-ip-updater service.
// It includes implementations for updating DNS records in AWS Route53 and other DNS providers.
package updater

import (
	"context"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Updater defines an interface for DNS record updates to allow for easy testing
// and potential support for multiple DNS providers. This interface can be implemented
// for different DNS services like AWS Route53, Cloudflare, etc.
type Updater interface {
	// Update performs the DNS record update operation.
	// It should update the specified DNS record with the new IP address.
	// Returns an error if the update operation fails.
	Update(context.Context) error
}

// AWSUpdater implements the Updater interface for AWS Route53 DNS service.
// It handles updating A records in a specified Route53 hosted zone.
type AWSUpdater struct {
	ZoneID    string // AWS Route53 hosted zone ID
	Subdomain string // Subdomain to update (e.g., "home.example.com")
	IP        string // New IP address to set for the subdomain
}

// Update performs a DNS record update in AWS Route53.
// It creates or updates an A record for the specified subdomain with the new IP address.
// The method uses AWS SDK v2 and requires proper AWS credentials to be configured.
//
// The update operation:
//   - Uses UPSERT action to create or update the A record
//   - Sets TTL to 60 seconds for quick propagation
//   - Adds a comment for tracking purposes
//   - Handles AWS API errors and returns them as Go errors
//
// Returns an error if:
//   - AWS credentials are invalid or missing
//   - The hosted zone ID is invalid
//   - The Route53 API call fails
//   - The subdomain format is invalid
func (awsupdater *AWSUpdater) Update(ctx context.Context) error {

	// Load AWS configuration from environment variables or IAM roles
	awscfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	// Create Route53 client with loaded configuration
	client := route53.NewFromConfig(awscfg)

	// Prepare the DNS record change request
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53types.ChangeBatch{
			Changes: []r53types.Change{
				{
					Action: r53types.ChangeActionUpsert, // Create or update the record
					ResourceRecordSet: &r53types.ResourceRecordSet{
						Name: aws.String(awsupdater.Subdomain), // Subdomain to update
						Type: r53types.RRTypeA,                 // A record type for IPv4 addresses
						TTL:  aws.Int64(60),                    // 60 second TTL for quick updates
						ResourceRecords: []r53types.ResourceRecord{
							{
								Value: aws.String(awsupdater.IP), // New IP address
							},
						},
					},
				},
			},
			Comment: aws.String("Updated by home-ip-updater"), // Tracking comment
		},
		HostedZoneId: aws.String(awsupdater.ZoneID), // Target hosted zone
	}

	// Execute the DNS record change
	_, errChange := client.ChangeResourceRecordSets(ctx, input)

	if errChange != nil {
		return errChange
	}

	return nil
}
