package config

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "View and get cli configuration",
}

func init() {
	cmd.RootCmd.AddCommand(configCmd)
}