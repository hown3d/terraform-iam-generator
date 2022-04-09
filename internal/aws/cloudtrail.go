package aws

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
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

const (
	defaultMaxRetries = 3
	apiCallEventType  = "AwsApiCall"
)

func NewCloudtrailService() (cloudTrailService, error) {
	cfg, err := newConfig(context.Background(), "eu-central-1")
	if err != nil {
		return cloudTrailService{}, fmt.Errorf("creating aws config: %w", err)
	}
	return cloudTrailService{
		client: newCloudtrail(cfg),
	}, nil
}

func (s cloudTrailService) GetIamActions(ctx context.Context, user string, startTime *time.Time, endTime *time.Time) (actions []IamAction, err error) {
	retries := 0
	nextToken := aws.String("")
	for nextToken != nil {
		resp, err := s.client.LookupEvents(ctx, &cloudtrail.LookupEventsInput{
			StartTime: startTime,
			EndTime:   endTime,
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
			return actions, nil
		}
		// aws cloudtrail can have a delay, so wait for the event to show off
		for _, event := range resp.Events {
			var s struct {
				UserAgent    string `json:"userAgent"`
				EventType    string `json:"eventType"`
				AWSRegion    string `json:"awsRegion"`
				UserIdentity struct {
					AccountID string `json:"accountId"`
				} `json:"userIdentity"`
			}
			json.Unmarshal([]byte(*event.CloudTrailEvent), &s)
			if hasTerraformUserAgent(s.UserAgent) && s.EventType == apiCallEventType {
				actions = append(actions, mapEventToAction(event, s.AWSRegion, s.UserIdentity.AccountID))
			}
		}
		if len(actions) == 0 && retries < defaultMaxRetries {
			retries++
			log.Printf("cloudtrail has delay uploading the events by up to 15 minutes, waiting for %v minutes before retrying", retries)
			time.Sleep(time.Duration(retries) * time.Minute)
			continue
		}
		nextToken = resp.NextToken
	}
	return actions, nil
}

type IamAction struct {
	Service   string
	APICall   string
	Resources []Arn
}
type Arn string

func mapEventToAction(e types.Event, awsRegion string, accountID string) IamAction {
	service := getServiceFromEventSource(*e.EventSource)
	resourceArns := make([]Arn, len(e.Resources))
	for _, r := range e.Resources {
		resourceArns = append(resourceArns, eventResourceArn(*r.ResourceName, getResourceTypeWithoutPrefix(*r.ResourceType), service, awsRegion, accountID))
	}
	return IamAction{
		APICall:   *e.EventName,
		Resources: resourceArns,
		Service:   service,
	}
}

func eventResourceArn(resourceName string, resourceType string, service string, region string, accountID string) Arn {
	return Arn(fmt.Sprintf("arn:aws:%s:%s:%s:%s/%s", service, region, accountID, resourceType, resourceName))
}

func getServiceFromEventSource(source string) string {
	return strings.Split(source, ".")[0]
}

// aws sets the resource Type by default to "AWS::<SERVICE>::<TYPE>".
// we only want the type in lowercase to use inside an arn
func getResourceTypeWithoutPrefix(resourceType string) string {
	return strings.ToLower(resourceType[strings.LastIndex(resourceType, "::")+2:])

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
