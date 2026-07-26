package route53

import (
	"context"
	"time"

	appconfig "github.com/a-castellano/home-ip-updater/internal/infra/config"
	aws "github.com/aws/aws-sdk-go-v2/aws"
	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

const TXTValue = "\"home-ip-updater-validation\""
const AWSRegion = "us-east-1"
const tracerName = "github.com/a-castellano/home-ip-updater/internal/infra/route53"

const awsRequestTimeout = 10 * time.Second

type Route53Updater struct {
	zoneID string // AWS Route53 hosted zone ID
	record string
	client *route53.Client
}

func (updater *Route53Updater) updateRoute53Record(ctx context.Context, recordType r53types.RRType, value string) error {
	// Prepare the DNS record change request
	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53types.ChangeBatch{
			Changes: []r53types.Change{
				{
					Action: r53types.ChangeActionUpsert, // Create or update the record
					ResourceRecordSet: &r53types.ResourceRecordSet{
						Name: aws.String(updater.record), // Subdomain to update
						Type: recordType,                 // A record type for IPv4 addresses
						TTL:  aws.Int64(60),              // 60 second TTL for quick updates
						ResourceRecords: []r53types.ResourceRecord{
							{
								Value: aws.String(value),
							},
						},
					},
				},
			},
			Comment: aws.String("Updated by home-ip-updater"), // Tracking comment
		},
		HostedZoneId: aws.String(updater.zoneID), // Target hosted zone
	}

	// Execute the DNS record change
	_, errChange := updater.client.ChangeResourceRecordSets(ctx, input)

	if errChange != nil {
		return errChange
	}

	return nil

}

func NewRoute53Updater(ctx context.Context, appConfig *appconfig.Config) (*Route53Updater, error) {

	var updater Route53Updater

	updater.zoneID = appConfig.AWSZoneID
	updater.record = appConfig.Subdomain

	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(AWSRegion),
		awsconfig.WithHTTPClient(awshttp.NewBuildableClient().WithTimeout(awsRequestTimeout)),
	)
	if err != nil {
		return nil, err
	}

	updater.client = route53.NewFromConfig(awscfg)

	//Check client updating TXT record
	updateErr := updater.updateRoute53Record(ctx, r53types.RRTypeTxt, TXTValue)

	if updateErr != nil {
		return nil, updateErr
	}

	return &updater, nil
}
