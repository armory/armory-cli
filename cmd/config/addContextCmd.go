package config

import (
	"fmt"
	"github.com/armory/armory-cli/internal/config"
	"github.com/spf13/cobra"
)

var AddContextCommand = &cobra.Command{
	Use:   "add-context",
	Short: "Add a context to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Add(cmd); err != nil {
			fmt.Printf("Unable to add context: %s\n", err.Error())
		}
	},
}

func init() {
	AddContextCommand.Flags().String(config.ParamAddClientId, "", "Armory cloud client-id")
	AddContextCommand.Flags().String(config.ParamAddSecret, "", "Armory cloud secret")
	AddContextCommand.Flags().String(config.ParamAddAudience, "", "Override Armory Cloud audience")
	AddContextCommand.Flags().String(config.ParamAddTokenIssuerUrl, "", "Override Armory Cloud token issuer endpoint")

	_ = AddContextCommand.MarkFlagRequired(config.ParamAddSecret)
	_ = AddContextCommand.MarkFlagRequired(config.ParamAddClientId)
	configBaseCommand.AddCommand(AddContextCommand)
}