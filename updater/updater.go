package updater

import (
	"context"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	route53 "github.com/aws/aws-sdk-go-v2/service/route53"
	r53types "github.com/aws/aws-sdk-go-v2/service/route53/types"
)

// Updater defines an interface in order to mock aws call inother libraries
type Updater interface {
	Update(context.Context) error
}

// AWSUpdater defines main struct that implements Updater interface
type AWSUpdater struct {
	ZoneID    string
	Subdomain string
	IP        string
}

// Update updates route53 record
func (awsupdater *AWSUpdater) Update(ctx context.Context) error {

	awscfg, err := awsconfig.LoadDefaultConfig(ctx)
	if err != nil {
		return err
	}

	client := route53.NewFromConfig(awscfg)

	input := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &r53types.ChangeBatch{
			Changes: []r53types.Change{
				{
					Action: r53types.ChangeActionUpsert,
					ResourceRecordSet: &r53types.ResourceRecordSet{
						Name: aws.String(awsupdater.Subdomain),
						Type: r53types.RRTypeA,
						TTL:  aws.Int64(60),
						ResourceRecords: []r53types.ResourceRecord{
							{
								Value: aws.String(awsupdater.IP),
							},
						},
					},
				},
			},
			Comment: aws.String("Updated by home-ip-updater"),
		},
		HostedZoneId: aws.String(awsupdater.ZoneID),
	}
	_, errChange := client.ChangeResourceRecordSets(ctx, input)

	if errChange != nil {
		return errChange
	}

	return nil
}
