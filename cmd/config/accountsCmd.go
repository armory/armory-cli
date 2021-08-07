package config

import (
	"context"
	"fmt"
	"github.com/armory/armory-cli/internal/config"
	"github.com/armory/armory-cli/internal/deng"
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
	"os"
	"time"
)

const (
	ParameterAccount  = "account"
	ParameterProvider = "provider"
)

var listAccountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "List available accounts",
	RunE: executeListAccountsCmd,
}

func init() {
	listAccountsCmd.Flags().String(config.ParamProvider, "kubernetes", "Provider")
	configCmd.AddCommand(listAccountsCmd)
}

func executeListAccountsCmd(cmd *cobra.Command, args []string) error {
	account, err := cmd.Flags().GetString(ParameterAccount)
	if err != nil {
		return err
	}
	provider, err := cmd.Flags().GetString(ParameterProvider)
	if err != nil {
		return err
	}
	ctx := context.TODO()
	client, err := helpers.MakeDeploymentClient(ctx, cmd)
	if err != nil {
		return err
	}

	r, err := client.GetApplications(ctx, &deng.GetAppRequest{
		Env: &deng.Environment{Provider: provider, Account: account},
	})
	if err != nil {
		return err
	}

	w := os.Stdout
	_, _ = fmt.Fprintf(w, "\nAPPLICATIONS\n")
	wt := ansiterm.NewTabWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprint(wt, "Name\tDeployments\tLast successful\tLast Failure\n")
	for _, c := range r.Apps {
		_, _ = fmt.Fprintf(wt, "%s\t%d\t%s\t%s\n", c.Name, c.Deployments, timeAsString(c.LastSuccessful), timeAsString(c.LastFailure))
	}
	_ = wt.Flush()
	return nil
}

func timeAsString(t *timestamp.Timestamp) string {
	tm := t.AsTime()
	if tm.IsZero() {
		return "-"
	}
	return tm.Local().Format(time.RFC822)
}