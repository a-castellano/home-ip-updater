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

	logger "github.com/a-castellano/go-services/infra/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
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

	ctx, span := otel.Tracer(tracerName).Start(ctx, "updateRoute53Record",
		trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	log := logger.FromContext(ctx).With("operation", "updateRoute53Record")
	log.InfoContext(ctx, "updating dns record using route53 updater", "recordType", recordType, "value", value)

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
		errorMessage := "an error hapended during route 53 record update"
		span.RecordError(errChange)
		span.SetStatus(codes.Error, errorMessage)
		log.ErrorContext(ctx, errorMessage, "error", errChange, "recordType", recordType, "value", value)
		return errChange
	}

	log.InfoContext(ctx, "dns record updated", "recordType", recordType, "value", value)
	return nil
}

func (updater *Route53Updater) UpdateRecord(ctx context.Context, value string) error {
	return updater.updateRoute53Record(ctx, r53types.RRTypeA, value)
}

func NewRoute53Updater(ctx context.Context, appConfig *appconfig.Config) (*Route53Updater, error) {

	var updater Route53Updater

	ctx, span := otel.Tracer(tracerName).Start(ctx, "NewRoute53Updater",
		trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	log := logger.FromContext(ctx).With("operation", "NewRoute53Updater")
	log.InfoContext(ctx, "setting up new route53 updater")

	updater.zoneID = appConfig.AWSZoneID
	updater.record = appConfig.Subdomain

	log.DebugContext(ctx, "loading AWS config")
	awscfg, err := awsconfig.LoadDefaultConfig(ctx,
		awsconfig.WithRegion(AWSRegion),
		awsconfig.WithHTTPClient(awshttp.NewBuildableClient().WithTimeout(awsRequestTimeout)),
	)
	if err != nil {
		errorMessage := "an error hapended during aws config generation"
		span.RecordError(err)
		span.SetStatus(codes.Error, errorMessage)
		log.ErrorContext(ctx, errorMessage, "error", err)
		return nil, err
	}

	log.DebugContext(ctx, "creatig route53 client")
	updater.client = route53.NewFromConfig(awscfg)

	log.DebugContext(ctx, "validating route53 client, updating TXT domain")
	//Check client updating TXT record
	updateErr := updater.updateRoute53Record(ctx, r53types.RRTypeTxt, TXTValue)

	if updateErr != nil {
		errorMessage := "cannot update TXT record using current configuration"
		// Status only: the error event and the log are already
		// recorded closest to the point of error
		span.SetStatus(codes.Error, errorMessage)
		log.ErrorContext(ctx, errorMessage, "error", errorMessage)
		return nil, updateErr
	}

	log.InfoContext(ctx, "route53 updater created")
	return &updater, nil
}
