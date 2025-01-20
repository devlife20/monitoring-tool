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
		ui.Show()
	},
}

func init() {

}
