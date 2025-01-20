package AWS

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
)

type CloudWatchLogsClient interface {
	GetLogEvents(ctx context.Context, params *cloudwatchlogs.GetLogEventsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.GetLogEventsOutput, error)
	DescribeLogGroups(ctx context.Context, params *cloudwatchlogs.DescribeLogGroupsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogGroupsOutput, error)
	DescribeLogStreams(ctx context.Context, params *cloudwatchlogs.DescribeLogStreamsInput, optFns ...func(*cloudwatchlogs.Options)) (*cloudwatchlogs.DescribeLogStreamsOutput, error)
}

func fetchAndPrintLogs(ctx context.Context, awsClient cloudwatchlogs.Client, groupName, streamName string, nextToken *string) (*string, error) {
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(groupName),
		LogStreamName: aws.String(streamName),
		StartFromHead: aws.Bool(true),
	}

	if nextToken != nil {
		input.NextToken = nextToken
	}

	resp, err := awsClient.GetLogEvents(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("error retrieving log events: %v", err)
	}

	for _, event := range resp.Events {
		fmt.Printf("[%v] %v\n",
			time.Unix(*event.Timestamp/1000, 0).Format(time.RFC3339),
			*event.Message)
	}

	return resp.NextForwardToken, nil
}

func TailCloudWatchLogs(groupName, streamName string, limit int32, region string) {
	ctx := context.Background()

	awsClient, err := AwsClient(ctx, region)
	if err != nil {
		log.Fatalf("Failed to connect to aws: %v", err)
	}

	// Get historical logs first
	fmt.Println("Fetching logs...")
	var nextToken *string
	token, err := fetchAndPrintLogs(ctx, awsClient, groupName, streamName, nextToken)
	if err != nil {
		log.Printf("Error fetching logs: %v", err)
	}
	nextToken = token

	fmt.Println("\nWatching for new logs...")
	fmt.Println("Press Ctrl+C to stop")

	// Then start tailing
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			token, err := fetchAndPrintLogs(ctx, awsClient, groupName, streamName, nextToken)
			if err != nil {
				log.Printf("Error fetching logs: %v", err)
				continue
			}
			nextToken = token
		}
	}
}
