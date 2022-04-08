package aws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudtrail"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func newConfig(ctx context.Context, region string) (aws.Config, error) {
	var cfg aws.Config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		return cfg, fmt.Errorf("loading default config: %w", err)
	}
	return cfg, nil
}

func newCloudtrail(cfg aws.Config) *cloudtrail.Client {
	return cloudtrail.NewFromConfig(cfg)
}

func newSTS(cfg aws.Config) *sts.Client {
	return sts.NewFromConfig(cfg)
}
