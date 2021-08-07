package deploy

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/spf13/cobra"
)

var deployCommand = &cobra.Command {
	Use:   "deploy",
	Short: "Initiate and manage deployments",
	Long:  `Initiate and manage deployments - see https://docs.armory.io/docs/deploy-engine`,
}

func init() {
	cmd.RootCmd.AddCommand(deployCommand)
}