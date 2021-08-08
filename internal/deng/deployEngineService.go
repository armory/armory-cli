package deng

import (
	"context"
	"github.com/armory-io/go-cloud-service/pkg/client"
	"github.com/armory/armory-cli/internal/config"
	"github.com/armory/armory-cli/internal/deng/protobuff"
	log "github.com/sirupsen/logrus"
	connection "google.golang.org/genproto/googleapis/firebase/fcm/connection/v1alpha1"
	"os"
	"os/signal"
)

type DeployEngineService struct {
	connection client.Connection
	deployEngineClient protobuff.DeploymentServiceClient
	ctx context.Context
}

var deployEngineService *DeployEngineService

func init() {
	conn, err := config.GetClientConnection(cmd)
	if err != nil {
		return nil, err
	}
	var cancelFunc context.CancelFunc
	ctx, cancelFunc := context.WithCancel(context.TODO())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	go func() {
		<-signalCh
		log.Debug("Signal interrupt received, stopping Deploy Engine command...")
		cancelFunc()
	}()

	// TODO can we connect lazily? maybe we don't for analytics?
	connection.Connect(ctx)
	deployEngineClient := protobuff.NewDeploymentServiceClient(connection.GetConn())
}

func GetApplications(provider string, account string) ([]*protobuff.AppSummary, error) {
	getAppResponse, err := deployEngineClient.GetApplications(ctx, &protobuff.GetAppRequest{
		Env: &protobuff.Environment{Provider: provider, Account: account},
	})
	if err != nil {
		return nil, err
	}
	return getAppResponse.Apps, nil
}

