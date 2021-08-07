package app

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/spf13/cobra"
)

const (
	ParameterAccount  = "account"
	ParameterProvider = "provider"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "Get and manage applications",
}

func init() {
	cmd.RootCmd.AddCommand(appCommand)
}