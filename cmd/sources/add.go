package sources

import (
	"github.com/devlife20/monitoring-tool/cmd/ui"
	"github.com/spf13/cobra"
)

var (
	filePath string
)

// AddCmd represents the add command
var AddCmd = &cobra.Command{
	Use:   "add",
	Short: "add a log source",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Run(filePath)
	},
}

func init() {

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	AddCmd.PersistentFlags().StringVarP(&filePath, "path", "p", "", "The url to ping")
	SourceCmd.MarkFlagRequired("path")
	SourceCmd.AddCommand(AddCmd)

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// addCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
