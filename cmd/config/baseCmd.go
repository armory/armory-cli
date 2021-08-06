package config

import "github.com/spf13/cobra"

var configBaseCommand = &cobra.Command{
	Use:   "config",
	Short: "View and get cli configuration",
}