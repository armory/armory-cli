package configCmd

import (
	"github.com/spf13/cobra"
)

var BaseCmd = &cobra.Command{
	Use:   "config",
	Short: "View and get cli configuration",
}

func init() {
	BaseCmd.AddCommand(addContextCommand)
	BaseCmd.AddCommand(listAccountsCmd)
	BaseCmd.AddCommand(viewCmd)
}