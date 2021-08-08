package deployCmd

import (
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/armory/armory-cli/internal/rollout"
	"github.com/spf13/cobra"
)

var abortCmd = &cobra.Command{
	Use:   "abort",
	Short: "Abort a deployment",
	Long: "After aborting a deployment, the resource will be put in the state where it can restarted with " +
		"armory deploy restart. This is different than rollback that reverts the deployment to the last known " +
		"good deployment.",
	RunE: func(c *cobra.Command, args []string) error {
		return helpers.ExecuteCancelable(c, rollout.Abort, args)
	},
}

func init() {
	abortCmd.Flags().String(rollout.ParamName, "", "name of deployment to abort")
	abortCmd.Flags().String(rollout.ParamType, "Deployment.apps", "type of deployment to abort")
}