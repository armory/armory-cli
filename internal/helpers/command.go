package helpers

import (
	"context"
	"github.com/armory/armory-cli/internal/config"
	"github.com/armory/armory-cli/internal/deng"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func MakeDeploymentClient(log *logrus.Logger, ctx context.Context, cmd *cobra.Command) (deng.DeploymentServiceClient, error) {
	conn, err := config.GetClientConnection(log, cmd)
	if err != nil {
		return nil, err
	}
	conn.Connect(ctx)
	return deng.NewDeploymentServiceClient(conn.GetConn()), nil
}