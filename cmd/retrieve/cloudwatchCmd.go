package retrieve

import (
	"fmt"
	"github.com/devlife20/monitoring-tool/CBL/AWS"
	"log"
	"time"

	"github.com/spf13/cobra"
)

var (
	logGroupName  string
	logStreamName string
	filterPattern string
	startTime     string
	endTime       string
	limit         int32
	region        string
)

var cloudCmd = &cobra.Command{
	Use:   "cloudwatch",
	Short: "Retrieve CloudWatch logs",
	Long: `Fetch logs from AWS CloudWatch with various filtering options.
   
Examples:
 monit cloudwatch -g my-log-group -s 'my-stream' -r eu-north-1
 monit cloudwatch -g my-log-group -p "ERROR" -r eu-north-1 -l 50 
 monit cloudwatch -g my-log-group -r eu-north-1 --start 2024-12-10T15:00:00Z --end 2024-12-10T16:00:00Z`,

	Run: func(cmd *cobra.Command, args []string) {
		// Convert start and end times to appropriate format
		if startTime != "" || endTime != "" {
			startTimeParsed, err := parseTime(startTime)
			if err != nil {
				log.Fatalf("Invalid start time: %v", err)
			}
			endTimeParsed, err := parseTime(endTime)
			if err != nil {
				log.Fatalf("Invalid end time: %v", err)
			}
			fmt.Printf("Start Time: %s\n", startTimeParsed)
			fmt.Printf("End Time: %s\n", endTimeParsed)
		}

		// Print the parameters for demonstration
		fmt.Printf("Fetching logs from AWS CloudWatch:\n")
		fmt.Printf("Log Group: %s\n", logGroupName)
		fmt.Printf("Log Stream: %s\n", logStreamName)
		fmt.Printf("Region: %s\n", region)
		fmt.Printf("Limit: %d\n", limit)

		AWS.GetCloudWatchLogs(logGroupName, logStreamName, filterPattern, limit, region)
	},
}

func init() {
	FetchLogCmd.AddCommand(cloudCmd)
	//  flags for the cloudCmd
	cloudCmd.Flags().StringVarP(&logGroupName, "group", "g", "", "Name of the CloudWatch log group (required)")
	cloudCmd.Flags().StringVarP(&logStreamName, "stream", "s", "", "Name of the CloudWatch log stream (optional)")
	cloudCmd.Flags().StringVarP(&filterPattern, "pattern", "p", "", "Pattern to filter log events (optional)")
	cloudCmd.Flags().StringVarP(&region, "region", "r", "", "AWS region for CloudWatch (required)")
	cloudCmd.Flags().StringVarP(&startTime, "start", "", "", "Start time for logs (RFC3339 format, e.g., 2024-12-10T15:00:00Z)")
	cloudCmd.Flags().StringVarP(&endTime, "end", "", "", "End time for logs (RFC3339 format, e.g., 2024-12-10T16:00:00Z)")
	cloudCmd.Flags().Int32VarP(&limit, "limit", "l", 100, "Max number of log events to fetch")

	// required flags
	// TODO: Set up a custom logger to redirect stderr to /var/log/monit/error.log.
	err := cloudCmd.MarkFlagRequired("group")
	if err != nil {
		return
	}
	err = cloudCmd.MarkFlagRequired("region")
	if err != nil {
		return
	}
}

// parseTime parses a time string in RFC3339 format.
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
