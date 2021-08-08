package deployCmd

import (
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/armory/armory-cli/internal/rollout"
	"github.com/spf13/cobra"
)

var restartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart a previously aborted deployment",
	RunE: func(c *cobra.Command, args []string) error {
		return helpers.ExecuteCancelable(c, rollout.Restart, args)
	},
}

func init() {
	restartCmd.Flags().String(rollout.ParamName, "", "name of deployment to restart")
	restartCmd.Flags().String(rollout.ParamType, "Deployment.apps", "type of deployment to restart")
}