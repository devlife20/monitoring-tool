/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/

package retrieve

import (
	"errors"
	"fmt"
	"github.com/devlife20/monitoring-tool/ELK"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
)

// elasticSearchCmd represents the elasticSearch command
var elasticSearchCmd = &cobra.Command{
	Use:   "elastic",
	Short: "Fetch logs from Elasticsearch",
	Long: `The "elastic" command allows you to retrieve logs stored in an Elasticsearch index. 
You can specify query parameters such as index name, search query, and time range 
to filter the logs effectively.

Examples:
  # Fetch all logs from the "app-logs" index
  fetch-log elastic --index app-logs

  # Fetch logs containing the term "ERROR" in the "app-logs" index
  fetch-log elastic --index app-logs --query "ERROR"

  # Fetch logs from a specific time range
  fetch-log elastic --index app-logs --start-time "2025-01-01T00:00:00" --end-time "2025-01-05T23:59:59"
`,
	PreRun: func(cmd *cobra.Command, args []string) {
		if err := viper.ReadInConfig(); err != nil {
			var configFileNotFoundError viper.ConfigFileNotFoundError
			if errors.As(err, &configFileNotFoundError) {
				// Configuration file not found
				log.Fatalf(`Configuration file not found.
Please run the command 'monit config init' to set up the configuration.`)
			}
		}

		// Check if elastic_api_key exists
		apiKey := viper.GetString("elastic_api_key")
		if apiKey == "" {
			log.Fatalf("Error: elastic_api_key not found in the configuration. Run 'fetch-log config' to set it.")
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		// Retrieve the flag values
		index, _ := cmd.Flags().GetString("index")
		query, _ := cmd.Flags().GetString("query")
		startTime, _ := cmd.Flags().GetString("start-time")
		endTime, _ := cmd.Flags().GetString("end-time")

		// Validate required flags
		if index == "" {
			log.Fatal("Error: --index flag is required")
		}

		// Simulate fetching logs
		fmt.Printf("Fetching logs from index: %s\n", index)
		if query != "" {
			fmt.Printf("Query: %s\n", query)
		}
		if startTime != "" || endTime != "" {
			fmt.Printf("Time Range: %s to %s\n", startTime, endTime)
		}

		err := ELK.FetchLogs(index)
		if err != nil {
			fmt.Println("=======", err)
			return
		}

		//ELK.UploadSyslogToElasticsearch("/var/log/syslog", "syslog-index")
	},
}

func init() {
	FetchLogCmd.AddCommand(elasticSearchCmd)

	// flags for the elastic command
	elasticSearchCmd.Flags().StringP("index", "i", "", "Name of the Elasticsearch index (required)")
	elasticSearchCmd.Flags().StringP("query", "q", "", "Search query to filter logs")
	elasticSearchCmd.Flags().String("start-time", "", "Start time for the logs (e.g., '2025-01-01T00:00:00')")
	elasticSearchCmd.Flags().String("end-time", "", "End time for the logs (e.g., '2025-01-05T23:59:59')")

	// --index as a required flag
	_ = elasticSearchCmd.MarkFlagRequired("index")
}
