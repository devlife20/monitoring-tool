package AWS

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs"
	_ "github.com/aws/aws-sdk-go-v2/service/cloudwatchlogs/types"
)

func GetCloudWatchLogs(groupName, streamName, region string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	awsClient, err := AwsClient(ctx, region)
	if err != nil {
		log.Fatalf("Failed to connect to aws: %v", err)
	}

	//  verify the log group exists
	_, err = awsClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(groupName),
	})
	if err != nil {
		log.Fatalf("Error describing log groups: %v", err)
	}

	// Then verify the log stream exists
	_, err = awsClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
		LogGroupName:        aws.String(groupName),
		LogStreamNamePrefix: aws.String(streamName),
	})
	if err != nil {
		log.Fatalf("Error describing log streams: %v", err)
	}

	// Retrieve log events
	input := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(groupName),
		LogStreamName: aws.String(streamName),
		Limit:         aws.Int32(int32(100)),
		StartFromHead: aws.Bool(true),
	}

	// Call the GetLogEvents API
	resp, err := awsClient.GetLogEvents(ctx, input)
	if err != nil {
		log.Fatalf("Error retrieving log events: %v", err)
	}

	//TODO  //implement a robust way to view the logs
	for _, event := range resp.Events {
		fmt.Printf("Timestamp: %v, Message: %v\n",
			time.Unix(*event.Timestamp/1000, 0),
			*event.Message)
	}
}
