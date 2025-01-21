package sources

import (
	"github.com/devlife20/monitoring-tool/cmd/ui"
	"github.com/spf13/cobra"
)

var (
	filePath string
)

// AddSourceCmd AddCmd represents the add command
var AddSourceCmd = &cobra.Command{
	Use:   "add-source",
	Short: "Add a new log source for monitoring",
	Long: `The "add-source" command allows you to register a new log source for real-time monitoring and analysis. 
This ensures logs from the specified source are included in your log management workflow.`,
	Run: func(cmd *cobra.Command, args []string) {
		ui.Show()
	},
}

func init() {

}
