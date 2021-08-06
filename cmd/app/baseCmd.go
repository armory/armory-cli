package app

import "github.com/spf13/cobra"

const (
	ParameterAccount  = "account"
	ParameterProvider = "provider"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "Get and manage applications",
}