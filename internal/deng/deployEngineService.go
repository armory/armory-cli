package deng

import (
	"context"
	"github.com/armory/armory-cli/internal/deng/protobuff"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"os"
	"os/signal"
)

type DeployEngineService struct {
	client     protobuff.DeploymentServiceClient
	ctx        context.Context
}

var deployEngineService *DeployEngineService

func GetDeployEngineInstance() *DeployEngineService {
	if deployEngineService == nil {
		ctx, cancel := context.WithCancel(context.TODO())

		signalCh := make(chan os.Signal, 1)
		signal.Notify(signalCh, os.Interrupt)

		go func() {
			<-signalCh
			log.Debug("Signal interrupt received, stopping Deploy Engine command...")
			cancel()
		}()

		var opts []grpc.DialOption
		deConnection, err := grpc.Dial("https://deploy-engin.stagin.cloud.armory.io:443", opts...)
		if err != nil {
			panic("Failed to create deploy engine connection")
		}
		deployEngineClient := protobuff.NewDeploymentServiceClient(deConnection)
		deployEngineService = &DeployEngineService{
			client: deployEngineClient,
			ctx: ctx,
		}
	}

	return deployEngineService
}

// GetApplications retrieves the applications for the given provider and account from Deploy Engine
func (deng *DeployEngineService) GetApplications(provider string, account string) ([]*AppSummary, error) {
	getAppResponse, err := deng.client.GetApplications(deng.ctx, &protobuff.GetAppRequest{
		Env: &protobuff.Environment{Provider: provider, Account: account},
	})
	if err != nil {
		return nil, err
	}

	var appSummaries []*AppSummary
	for i := range getAppResponse.Apps {
		var app = getAppResponse.Apps[i]
		appSummaries = append(appSummaries, &AppSummary{
			Name: app.Name,
			Deployments: app.Deployments,
			LastSuccessful: app.LastSuccessful.AsTime(),
			LastFailure: app.LastFailure.AsTime(),
		})
	}
	return appSummaries, nil
}

// GetAccounts gets a list of account summaries for a given provider
func (deng *DeployEngineService) GetAccounts(provider string) ([]*AccountSummary, error) {
	getAccountResponse, err := deng.client.GetAccounts(deng.ctx, &protobuff.GetAccountRequest{Provider: provider})
	if err != nil {
		return nil, err
	}
	var accounts []*AccountSummary
	for i := range getAccountResponse.Accounts {
		var account = getAccountResponse.Accounts[i]
		accounts = append(accounts, &AccountSummary{
			Provider: account.Provider,
			Name:     account.Account,
		})
	}
	return accounts, nil
}

// ResumeDeployment Resumes a partial or complete deployment\n
// deploymentId The deployment id
// atomicDeploymentPartName The name of the atomic deployment part of the deployment to resume
// atomicDeploymentPartType The type of the atomic deployment part of the deployment to resume
func (deng *DeployEngineService) ResumeDeployment(deploymentId string, atomicDeploymentPartName string, atomicDeploymentPartType string) error {
	req := createRolloutRequest(deploymentId, atomicDeploymentPartName, atomicDeploymentPartType)
	res, err := deng.client.Resume(deng.ctx, req)
	if err != nil {
		return nil
	}
	log.Info(res.Message)
	return nil
}

// RestartDeployment Restarts a partial or complete deployment
// deploymentId The deployment id
// atomicDeploymentPartName The name of the atomic deployment part of the deployment to resume
// atomicDeploymentPartType The type of the atomic deployment part of the deployment to resume
func (deng *DeployEngineService) RestartDeployment(deploymentId string, atomicDeploymentPartName string, atomicDeploymentPartType string) error {
	req := createRolloutRequest(deploymentId, atomicDeploymentPartName, atomicDeploymentPartType)
	res, err := deng.client.Restart(deng.ctx, req)
	if err != nil {
		return nil
	}
	log.Info(res.Message)
	return nil
}

// AbortDeployment Aborts a partial or complete deployment
// deploymentId The deployment id
// atomicDeploymentPartName The name of the atomic deployment part of the deployment to resume
// atomicDeploymentPartType The type of the atomic deployment part of the deployment to resume
func (deng *DeployEngineService) AbortDeployment(deploymentId string, atomicDeploymentPartName string, atomicDeploymentPartType string) error {
	req := createRolloutRequest(deploymentId, atomicDeploymentPartName, atomicDeploymentPartType)
	res, err := deng.client.Resume(deng.ctx, req)
	if err != nil {
		return nil
	}
	log.Info(res.Message)
	return nil
}

// StartDeployment Starts a deployment and returns the deployment ID.
func (deng *DeployEngineService) StartDeployment(deploymentConfiguration *DeploymentConfiguration) (string, error) {
	request := createDeploymentRequestFromConfiguration(deploymentConfiguration)
	deploymentDescriptor, err := deng.client.Start(deng.ctx, request)
	if err != nil {
		return "", err
	}

	deploymentId := deploymentDescriptor.Id
	return deploymentId, nil
}

// createRolloutRequest Creates the Rollout requests
func createRolloutRequest(deploymentId string, atomicDeploymentPartName string, atomicDeploymentPartType string) *protobuff.RolloutRequest {
	var req *protobuff.RolloutRequest
	if atomicDeploymentPartName == "" {
		req = &protobuff.RolloutRequest{All: true, DeploymentId: deploymentId}
	} else {
		req = &protobuff.RolloutRequest{
			Name: atomicDeploymentPartName,
			Type: atomicDeploymentPartType,
			DeploymentId: deploymentId,
		}
	}
	return req
}


func createDeploymentRequestFromConfiguration(deploymentConfiguration *DeploymentConfiguration) *protobuff.Deployment {

}