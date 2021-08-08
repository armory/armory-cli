package appCmd

import (
	"github.com/spf13/cobra"
)

const (
	ParameterAccount  = "account"
	ParameterProvider = "provider"
)

var BaseCmd = &cobra.Command{
	Use:   "app",
	Short: "Get and manage applications",
}

func init() {
	BaseCmd.AddCommand(listApplicationsCommand)
}