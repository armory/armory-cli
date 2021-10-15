package accounts

import (
	"github.com/armory/armory-cli/cmd"
	"github.com/spf13/cobra"
)

const (
	accountsShort   = "Obtain list of existing accounts."
	accountsLong    = ""
	accountsExample    = ""
)

type accountsOptions struct {
	*cmd.RootOptions
	providerName string
}

func NewAccountsCmd(rootOptions *cmd.RootOptions) *cobra.Command {
	options := &accountsOptions{
		RootOptions: rootOptions,
	}
	command := &cobra.Command{
		Use:     "accounts",
		Aliases: []string{"accounts"},
		Short:   accountsShort,
		Long:    accountsLong,
		Example: accountsExample,
	}
	// create subcommands
	command.AddCommand(NewListCmd(options))

	return command
}

