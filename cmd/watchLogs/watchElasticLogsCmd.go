package watchLogs

import (
	"github.com/devlife20/monitoring-tool/ELK"
	"github.com/spf13/cobra"
	"time"
)

var (
	// Command flags
	elasticURL   string
	indexName    string
	queryString  string
	interval     time.Duration
	fromDate     string
	toDate       string
	maxResults   int
	sortField    string
	sortOrder    string
	outputFormat string
	basicAuth    bool
	username     string
	password     string
)

var watchElastic = &cobra.Command{
	Use:   "elastic",
	Short: "Watch logs from elasticsearch in real time or on a scheduled basis",
	Long: `
### Authentication Methods:
- Supports API key authentication. Ensure credentials are configured using the 'monit config elastic' command.
    `,
	Example: `
    # Watch logs in real-time from a specific index
    monit watch elastic --url localhost:9200 --index nginx-logs --realtime

    # Query logs from the last hour with custom interval
    monit watch elastic --url localhost:9200 --index apache-logs --interval 5m --from-date "1h ago"`,

	Run: func(cmd *cobra.Command, args []string) {
		ELK.WatchElasticLogs(indexName)
	},
}

func init() {
	WatchCmd.AddCommand(watchElastic)
	// Required flags
	watchElastic.Flags().StringVar(&elasticURL, "url", "", "Elasticsearch server URL (required)")
	watchElastic.Flags().StringVar(&indexName, "index", "", "Name of the Elasticsearch index to query (required)")
	watchElastic.Flags().DurationVar(&interval, "interval", 0, "Polling interval for log retrieval (e.g., 30s, 5m, 1h)")
	watchElastic.Flags().StringVar(&fromDate, "from-date", "", "Start date for log retrieval (format: YYYY-MM-DD or relative time like '1h ago')")
	// Mark required flags
	//watchElastic.MarkFlagRequired("url")
	watchElastic.MarkFlagRequired("index")
}
