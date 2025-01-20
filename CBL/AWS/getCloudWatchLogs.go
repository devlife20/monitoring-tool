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

func GetCloudWatchLogs(groupName, streamName, filterPattern string, limit int32, region string) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	awsClient, err := AwsClient(ctx, region)
	if err != nil {
		log.Fatalf("Failed to connect to aws: %v", err)
	}

	// Verify the log group exists
	_, err = awsClient.DescribeLogGroups(ctx, &cloudwatchlogs.DescribeLogGroupsInput{
		LogGroupNamePrefix: aws.String(groupName),
	})
	if err != nil {
		log.Fatalf("Error describing log groups: %v", err)
	}

	if filterPattern != "" {
		// Use FilterLogEvents when a filter pattern is provided
		filterInput := &cloudwatchlogs.FilterLogEventsInput{
			LogGroupName:  aws.String(groupName),
			FilterPattern: aws.String(filterPattern),
			Limit:         aws.Int32(limit),
		}

		// Add stream name to filter if provided
		if streamName != "" {
			filterInput.LogStreamNames = []string{streamName}
		}

		// Get filtered logs
		output, err := awsClient.FilterLogEvents(ctx, filterInput)
		if err != nil {
			log.Fatalf("Error filtering log events: %v", err)
		}

		for _, event := range output.Events {
			fmt.Printf("Timestamp: %v, Stream: %v, Message: %v\n",
				time.Unix(*event.Timestamp/1000, 0),
				*event.LogStreamName,
				*event.Message)
		}

		// Handle pagination if necessary
		for output.NextToken != nil {
			filterInput.NextToken = output.NextToken
			output, err = awsClient.FilterLogEvents(ctx, filterInput)
			if err != nil {
				log.Fatalf("Error filtering log events: %v", err)
			}

			for _, event := range output.Events {
				fmt.Printf("Timestamp: %v, Stream: %v, Message: %v\n",
					time.Unix(*event.Timestamp/1000, 0),
					*event.LogStreamName,
					*event.Message)
			}
		}
	} else if streamName != "" {
		// If no filter pattern but stream name is provided, use GetLogEvents
		// Verify the stream exists first
		_, err = awsClient.DescribeLogStreams(ctx, &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName:        aws.String(groupName),
			LogStreamNamePrefix: aws.String(streamName),
		})
		if err != nil {
			log.Fatalf("Error describing log streams: %v", err)
		}

		input := &cloudwatchlogs.GetLogEventsInput{
			LogGroupName:  aws.String(groupName),
			LogStreamName: aws.String(streamName),
			Limit:         aws.Int32(limit),
			StartFromHead: aws.Bool(true),
		}

		resp, err := awsClient.GetLogEvents(ctx, input)
		if err != nil {
			log.Fatalf("Error retrieving log events: %v", err)
		}

		for _, event := range resp.Events {
			fmt.Printf("Timestamp: %v, Message: %v\n",
				time.Unix(*event.Timestamp/1000, 0),
				*event.Message)
		}
	} else {
		// If neither filter pattern nor stream name is provided, list recent logs from all streams
		input := &cloudwatchlogs.DescribeLogStreamsInput{
			LogGroupName: aws.String(groupName),
			OrderBy:      "LastEventTime",
			Descending:   aws.Bool(true),
			Limit:        aws.Int32(1), // Get most recent stream
		}

		streams, err := awsClient.DescribeLogStreams(ctx, input)
		if err != nil {
			log.Fatalf("Error describing log streams: %v", err)
		}

		if len(streams.LogStreams) > 0 {
			mostRecentStream := streams.LogStreams[0]
			eventsInput := &cloudwatchlogs.GetLogEventsInput{
				LogGroupName:  aws.String(groupName),
				LogStreamName: mostRecentStream.LogStreamName,
				Limit:         aws.Int32(limit),
				StartFromHead: aws.Bool(true),
			}

			resp, err := awsClient.GetLogEvents(ctx, eventsInput)
			if err != nil {
				log.Fatalf("Error retrieving log events: %v", err)
			}

			for _, event := range resp.Events {
				fmt.Printf("Timestamp: %v, Message: %v\n",
					time.Unix(*event.Timestamp/1000, 0),
					*event.Message)
			}
		}
	}
}
