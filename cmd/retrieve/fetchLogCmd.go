/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/

package retrieve

import (
	"github.com/spf13/cobra"
)

// FetchLogCmd represents the "fetch" command.
// All commands in the "retrieve" package are subcommands of the "fetch" command.
// When the "fetch" command is called without any arguments, it displays the help message
var FetchLogCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Retrieve logs from various sources such as cloud, local storage, or ELK stack",
	Long: `The "fetch" command allows you to retrieve logs from different sources including:
  - Cloud providers (e.g., AWS, GCP, Azure)
  - Local directories or servers
  - Elasticsearch, Logstash, and Kibana (ELK stack)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			_ = cmd.Help()
		}
	},
}
