package deployCmd

import (
	"github.com/spf13/cobra"
)

var BaseCmd = &cobra.Command {
	Use:   "deploy",
	Short: "Initiate and manage deployments",
}

func init() {
	// Add deploy sub commands
	BaseCmd.AddCommand(abortCmd)
	BaseCmd.AddCommand(restartCmd)
	BaseCmd.AddCommand(resumeCmd)
	BaseCmd.AddCommand(startCmd)
	BaseCmd.AddCommand(statusCmd)
}