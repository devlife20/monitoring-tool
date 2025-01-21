package cmd

import (
	"github.com/devlife20/monitoring-tool/cmd/config"
	"github.com/devlife20/monitoring-tool/cmd/retrieve"
	"github.com/devlife20/monitoring-tool/cmd/sources"
	"github.com/devlife20/monitoring-tool/cmd/watchLogs"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "monit",
	Short: "monit is a simple CLI tool to manage and monitor logs from various sources",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//Run: func(cmd *cobra.Command, args []string) {
	//	fmt.Println("Hello world welcome to cobra")
	//},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.AddCommand(sources.AddSourceCmd)
	rootCmd.AddCommand(retrieve.FetchLogCmd)
	rootCmd.AddCommand(config.MonitConfig)
	rootCmd.AddCommand(watchLogs.WatchCmd)

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.monitoring-tool.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
