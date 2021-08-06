package config

import (
	"github.com/armory/armory-cli/internal/config"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View CLI configuration",
	Run: func(cmd *cobra.Command, args []string) {
		config.View()
	},
}

func init() {
	configBaseCommand.AddCommand(viewCmd)
}