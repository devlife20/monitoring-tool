package AWS

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

func AwsClient(ctx context.Context, region string) (cloudwatchlogs.Client, error) {
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
		config.WithRetryMode(aws.RetryModeStandard),
		config.WithRetryMaxAttempts(3),
	)
	if err != nil {
		return cloudwatchlogs.Client{}, fmt.Errorf("unable to load SDK config: %w", err)
	}
	client := cloudwatchlogs.NewFromConfig(cfg)

	return *client, nil

}
