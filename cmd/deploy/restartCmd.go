package deploy

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/armory/armory-cli/internal/rollout"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a previously aborted deployment",
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.ExecuteCancelable(c, rollout.Restart, args)
	},
}

func init() {
	getCmd.Flags().String(rollout.ParamName, "", "name of deployment to restart")
	getCmd.Flags().String(rollout.ParamType, "Deployment.apps", "type of deployment to restart")
	deployCommand.AddCommand(restartCmd)
}