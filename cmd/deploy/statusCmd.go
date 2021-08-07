package deploy

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/armory/armory-cli/internal/status"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "status",
	Short: "Get deployment information",
	Long:  `Get deployment information [deployment ID]`,
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.ExecuteCancelable(c, status.Execute, args)
	},
}

func init() {
	getCmd.Flags().BoolP(status.ParameterWatch, "w", false, "watch changes")
	getCmd.Flags().Bool(status.ParameterShowEvents, false, "show events")
	deployCommand.AddCommand(getCmd)
}