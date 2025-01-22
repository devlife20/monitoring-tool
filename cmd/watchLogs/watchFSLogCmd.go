/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package watchLogs

import (
	"github.com/devlife20/monitoring-tool/LFS/linux"
	"github.com/spf13/cobra"
	"log"
)

var (
	path       string
	logDirPath string
)

// watchCmd represents the watch command
var watchFS = &cobra.Command{
	Use:   "local",
	Short: "Continuously monitor and stream logs from a file",
	Long: `Monitor and continuously stream logs from a specified file. with optional filtering based on a specified pattern.

Examples:
  monit watch --path /var/log/myapp.log --pattern "ERROR"
  monit watch --path /var/log/myapp.log --pattern "INFO"
`,
	Run: func(cmd *cobra.Command, args []string) {
		err := linux.TailLogsFromFile(path, filterPattern)
		if err != nil {
			log.Fatalf("failed to watch file %s", err)
		}
	},
}

func init() {
	WatchCmd.AddCommand(watchFS)
	watchFS.Flags().StringVarP(&path, "path", "", "", "Path to the log file (required)")
	watchFS.Flags().StringVarP(&filterPattern, "pattern", "p", "", "Pattern to filter log lines (optional)")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// watchCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// watchCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
