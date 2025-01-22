package retrieve

import (
	"github.com/devlife20/monitoring-tool/LFS/linux"
	"github.com/spf13/cobra"
	"log"
)

var (
	path       string
	logDirPath string
)

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Retrieve logs from file system",
	Long: `Fetch logs from file system with various filtering options.

Examples:
 monit fetch local --path /var/log/myapp.log --pattern "ERROR"
 monit fetch local --path /var/log/myapp.log --pattern "INFO"`,

	Run: func(cmd *cobra.Command, args []string) {

		if path == "" {
			log.Fatal(" --path must be specified")
		}
		if path != "" {
			linux.FetchLogsFromFile(path, filterPattern)
		}
	},
}

func init() {
	FetchLogCmd.AddCommand(localCmd)

	// Flags for the localCmd
	localCmd.Flags().StringVarP(&path, "path", "", "", "Path to the log file (required)")
	localCmd.Flags().StringVarP(&filterPattern, "pattern", "p", "", "Pattern to filter log lines (optional)")

	// Mark --path flag as required
	localCmd.MarkFlagRequired("path")
}

// fetchLogsFromFile reads logs from a file and applies filters.
