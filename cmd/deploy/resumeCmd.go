package deploy

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/armory/armory-cli/internal/rollout"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume a paused deployment",
	RunE: func(c *cobra.Command, args []string) error {
		return cmd.ExecuteCancelable(c, rollout.Resume, args)
	},
}

func init() {
	getCmd.Flags().String(rollout.ParamName, "", "name of deployment to resume")
	getCmd.Flags().String(rollout.ParamType, "Deployment.apps", "type of deployment to resume")
	deployCommand.AddCommand(resumeCmd)
}