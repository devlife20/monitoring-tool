package watchLogs

import (
	"github.com/devlife20/monitoring-tool/CBL/AWS"
	"github.com/spf13/cobra"
	"time"
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

// WatchCmd represents the watch command
var watchCloudwatch = &cobra.Command{
	Use:   "cloudwatch",
	Short: "Watch logs from cloudwatch in real time or on a scheduled basis",
	Long: `aggregate logs from cloudwatch in real time or on a scheduled basis

	Examples:
	monit watch cloudwatch -g my-log-group -s 'my-stream' -r eu-north-1`,
	Run: func(cmd *cobra.Command, args []string) {
		AWS.TailCloudWatchLogs(logGroupName, logStreamName, region)
	},
}

func init() {
	WatchCmd.AddCommand(watchCloudwatch)
	watchCloudwatch.Flags().StringVarP(&logGroupName, "group", "g", "", "Name of the CloudWatch log group (required)")
	watchCloudwatch.Flags().StringVarP(&logStreamName, "stream", "s", "", "Name of the CloudWatch log stream (optional)")
	watchCloudwatch.Flags().StringVarP(&filterPattern, "pattern", "p", "", "Pattern to filter log events (optional)")
	watchCloudwatch.Flags().StringVarP(&region, "region", "r", "", "AWS region for CloudWatch (required)")
	watchCloudwatch.Flags().StringVarP(&startTime, "start", "", "", "Start time for logs (RFC3339 format, e.g., 2024-12-10T15:00:00Z)")
	watchCloudwatch.Flags().StringVarP(&endTime, "end", "", "", "End time for logs (RFC3339 format, e.g., 2024-12-10T16:00:00Z)")
	watchCloudwatch.Flags().Int32VarP(&limit, "limit", "l", 100, "Max number of log events to fetch")

	// required flags
	err := watchCloudwatch.MarkFlagRequired("group")
	if err != nil {
		return
	}
	err = watchCloudwatch.MarkFlagRequired("region")
	if err != nil {
		return
	}
}

// parseTime parses a time string in RFC3339 format.
func parseTime(timeStr string) (time.Time, error) {
	return time.Parse(time.RFC3339, timeStr)
}
