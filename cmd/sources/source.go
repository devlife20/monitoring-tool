package sources

import (
	"github.com/spf13/cobra"
)

// SourceCmd represents the source command
var SourceCmd = &cobra.Command{
	Use:   "source",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		err := cmd.Help()
		if err != nil {
			cobra.CheckErr(err)
		}

	},
}

func init() {
	SourceCmd.AddCommand(AddCmd)
}
