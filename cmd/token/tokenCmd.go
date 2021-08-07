package token

import (
	"context"
	"github.com/armory/armory-cli/cmd"
	"github.com/armory/armory-cli/internal/config"
	"github.com/spf13/cobra"
)

// TODO this command should be deleted probably
var tokenCommand = &cobra.Command{
	Use:   "token",
	Short: "Obtain a token from configured provider",
	RunE: func(cmd *cobra.Command, args []string) error {
		return config.GetToken(context.TODO(), cmd)
	},
}

func init() {
	cmd.RootCmd.AddCommand(tokenCommand)
}