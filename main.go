package main

import (
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	awssession "github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
)

func main() {

	zoneId := os.Getenv("ZONE_ID")
	subdomain := os.Getenv("SUBDOMAIN")

	sess, err := awssession.NewSession()
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	svc := route53.New(sess)

	fmt.Println("svc", svc)

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action: aws.String("UPSERT"),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(subdomain),
						ResourceRecords: []*route53.ResourceRecord{
							{
								Value: aws.String("192.0.2.44"),
							},
						},
						TTL:  aws.Int64(60),
						Type: aws.String("A"),
					},
				},
			},
			Comment: aws.String("Updated by ot"),
		},
		HostedZoneId: aws.String(zoneId),
	}

	resp, err := svc.ChangeResourceRecordSets(params)

	if err != nil {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		fmt.Println(err.Error())
		return
	}

	// Pretty-print the response data.
	fmt.Println("Change Response:")
	fmt.Println(resp)
}
