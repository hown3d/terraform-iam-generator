package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail/types"
	awsbase "github.com/hashicorp/aws-sdk-go-base/v2"
)

type CloudTrailClient interface {
	LookupEvents(ctx context.Context, params *cloudtrail.LookupEventsInput, optFns ...func(*cloudtrail.Options)) (*cloudtrail.LookupEventsOutput, error)
}

type cloudTrailService struct {
	client CloudTrailClient
}

func NewEventService() (cloudTrailService, error) {
	cfg, err := newConfig(context.Background(), "eu-central-1")
	if err != nil {
		return cloudTrailService{}, fmt.Errorf("creating aws config: %w", err)
	}
	return cloudTrailService{
		client: newCloudtrail(cfg),
	}, nil
}

func (s cloudTrailService) GetEventsOfUser(ctx context.Context, user string, startTime *time.Time, endTime *time.Time) (events []types.Event, err error) {
	nextToken := aws.String("")
	for nextToken != nil {
		resp, err := s.client.LookupEvents(ctx, &cloudtrail.LookupEventsInput{
			StartTime: startTime,
			EndTime:   endTime,
			//EventCategory: "managment",
			LookupAttributes: []types.LookupAttribute{
				{
					AttributeKey:   types.LookupAttributeKeyUsername,
					AttributeValue: aws.String(user),
				},
			},
		})
		if err != nil {
			return nil, fmt.Errorf("looking up events in cloudtrail: %w", err)
		}
		if resp == nil {
			return events, nil
		}
		for _, event := range resp.Events {
			var s struct {
				UserAgent string `json:"userAgent"`
			}
			json.Unmarshal([]byte(*event.CloudTrailEvent), &s)
			if hasTerraformUserAgent(s.UserAgent) {
				events = append(events, event)
			}
		}
		nextToken = resp.NextToken
	}
	return events, nil
}

type myEvent struct {
	types.Event
}

func hasTerraformUserAgent(userAgent string) bool {
	// https://github.com/hashicorp/terraform-provider-aws/blob/3ce53331abc4bac9bac4eca1cf477fe74bba5bb9/internal/conns/conns.go#L1339
	// taken from function StdUserAgentProducts
	apnInfo := awsbase.APNInfo{
		PartnerName: "HashiCorp",
		Products: []awsbase.UserAgentProduct{
			{
				Name: "Terraform",
			},
		},
	}
	defaultUserAgentString := apnInfo.BuildUserAgentString()
	// trim off any "[" or "]" since the userAgent string in aws has them
	userAgent = strings.Trim(userAgent, "[]")
	if strings.HasPrefix(userAgent, defaultUserAgentString) {
		return true
	}
	return false
}
