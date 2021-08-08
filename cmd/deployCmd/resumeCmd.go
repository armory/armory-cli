package deployCmd

import (
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/armory/armory-cli/internal/rollout"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused deployment",
	RunE: func(c *cobra.Command, args []string) error {
		return helpers.ExecuteCancelable(c, rollout.Resume, args)
	},
}

func init() {
	resumeCmd.Flags().String(rollout.ParamName, "", "name of deployment to resume")
	resumeCmd.Flags().String(rollout.ParamType, "Deployment.apps", "type of deployment to resume")
}