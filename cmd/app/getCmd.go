package app

import (
	"context"
	"fmt"
	"github.com/armory/armory-cli/internal/deng"
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/armory/armory-cli/internal/utils"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/juju/ansiterm"
	"github.com/spf13/cobra"
	"os"
	"time"
)

var listApplicationsCommand = &cobra.Command{
	Use:   "list",
	Short: "Lists application information",
	Long:  `Retrieve known applications.`,
	RunE: execute,
}

func init() {
	listApplicationsCommand.Flags().String(ParameterProvider, "kubernetes", "provider")
	listApplicationsCommand.Flags().String(ParameterAccount, "", "account name")
	appCommand.AddCommand(listApplicationsCommand)
}

func execute(cmd *cobra.Command, args []string) error {
	ctx := context.TODO()
	logger := utils.GetLogger()

	client, err := helpers.MakeDeploymentClient(logger, ctx, cmd)
	if err != nil {
		return err
	}

	envName, err := cmd.Flags().GetString(ParameterAccount)
	if err != nil {
		return err
	}
	envProvider, err := cmd.Flags().GetString(ParameterProvider)
	if err != nil {
		return err
	}

	r, err := client.GetApplications(ctx, &deng.GetAppRequest{
		Env: &deng.Environment{Provider: envProvider, Account: envName},
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
