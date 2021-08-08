package app

import (
	"fmt"
	"github.com/armory/armory-cli/internal/deng"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/juju/ansiterm"
	"os"
	"time"
)

func PrintApplications(provider string, account string) error {
	apps, err := deng.GetApplications(provider, account)
	if err != nil {
		return err
	}

	w := os.Stdout
	_, _ = fmt.Fprintf(w, "\nAPPLICATIONS\n")
	wt := ansiterm.NewTabWriter(w, 0, 0, 2, ' ', 0)
	_, _ = fmt.Fprint(wt, "Name\tDeployments\tLast successful\tLast Failure\n")
	for _, c := range apps {
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