package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type stsClient interface {
	GetCallerIdentity(ctx context.Context, params *sts.GetCallerIdentityInput, optFns ...func(*sts.Options)) (*sts.GetCallerIdentityOutput, error)
}

func NewStsService() (stsService, error) {
	cfg, err := newConfig(context.Background(), "eu-central-1")
	if err != nil {
		return stsService{}, fmt.Errorf("creating aws config: %w", err)
	}
	return stsService{
		client: newSTS(cfg),
	}, nil
}

type stsService struct {
	client stsClient
}

func (s stsService) GetCallerIdentity(ctx context.Context) (string, error) {
	resp, err := s.client.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
	if err != nil {
		return "", fmt.Errorf("getting caller identity: %w", err)
	}
	return *resp.Arn, nil
}
