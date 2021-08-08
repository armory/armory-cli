package configCmd

import (
	"fmt"
	"github.com/armory/armory-cli/internal/config"
	"github.com/spf13/cobra"
)

var addContextCommand = &cobra.Command{
	Use:   "add-context",
	Short: "Add a context to the configuration",
	Run: func(cmd *cobra.Command, args []string) {
		if err := config.Add(cmd); err != nil {
			fmt.Printf("Unable to add context: %s\n", err.Error())
		}
	},
}

func init() {
	addContextCommand.Flags().String(config.ParamAddClientId, "", "Armory cloud client-id")
	addContextCommand.Flags().String(config.ParamAddSecret, "", "Armory cloud secret")
	addContextCommand.Flags().String(config.ParamAddAudience, "", "Override Armory Cloud audience")
	addContextCommand.Flags().String(config.ParamAddTokenIssuerUrl, "", "Override Armory Cloud token issuer endpoint")

	_ = addContextCommand.MarkFlagRequired(config.ParamAddSecret)
	_ = addContextCommand.MarkFlagRequired(config.ParamAddClientId)
}