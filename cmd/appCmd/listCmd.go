package appCmd

import (
	"github.com/armory/armory-cli/internal/app"
	"github.com/spf13/cobra"
)

var listApplicationsCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists application information",
	Long:  `Retrieve known applications.`,
	RunE: executeListApplicationsCommand,
}

func init() {
	listApplicationsCommand.Flags().String(ParameterProvider, "kubernetes", "provider")
	listApplicationsCommand.Flags().String(ParameterAccount, "", "account name")
	BaseCmd.AddCommand(listApplicationsCommand)
}

func executeListApplicationsCommand(cmd *cobra.Command, args []string) error {
	account, err := cmd.Flags().GetString(ParameterAccount)
	if err != nil {
		return err
	}
	provider, err := cmd.Flags().GetString(ParameterProvider)
	if err != nil {
		return err
	}
	return app.PrintApplications(account, provider)
}