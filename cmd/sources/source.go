package sources

import (
	"github.com/spf13/cobra"
)

// SourceCmd represents the source command
var SourceCmd = &cobra.Command{
	Use:   "source",
	Short: "",
	Long: `The source command group handles the setup and management of log sources, 
allowing you to configure, list, test, and remove sources for log retrieval`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}
