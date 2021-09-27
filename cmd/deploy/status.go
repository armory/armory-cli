package deploy

import (
	"fmt"
	deploy "github.com/armory-io/deploy-engine/deploy/client"
	"github.com/armory/armory-cli/cmd"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

const (
	deployStatusShort   = "Watch deployment on Armory Cloud"
	deployStatusLong    = "Watch deployment on Armory Cloud"
	deployStatusExample = "armory deploy status [options]"
	cloudConsoleBaseUrl = "https://console.cloud.armory.io"
	cloudConsoleStagingBaseUrl = "https://console.staging.cloud.armory.io"//deployments/f6567ec2-012f-4dd1-ba89-f35530789d1a
)

type StatusOptions struct {
	DeploymentId   string
}

func NewDeployStatusCmd(deployOptions *cmd.RootOptions) *cobra.Command {
	statusOptions := &StatusOptions{}
	cmd := &cobra.Command{
		Use:     "status -deploymentId [deploymentId]",
		Aliases: []string{"status"},
		Short:   deployStatusShort,
		Long:    deployStatusLong,
		Example: deployStatusExample,
		RunE: func(cmd *cobra.Command, args []string) error {
			return status(cmd, deployOptions, statusOptions)
		},
	}
	cmd.PersistentFlags().StringVarP(&statusOptions.DeploymentId, "deploymentId", "i", "", "The id of an existing deployment (required)")
	cmd.MarkPersistentFlagRequired("deploymentId")
	return cmd
}

func status(cmd *cobra.Command, deployOptions *cmd.RootOptions, statusOptions *StatusOptions) error {
	req := deployOptions.DeployClient.DeploymentServiceApi.DeploymentServiceStatus(deployOptions.DeployClient.Context, statusOptions.DeploymentId)
	deployResp, response, err := req.Execute()
	if err != nil && response.StatusCode >= 300 {
		openAPIErr := err.(deploy.GenericOpenAPIError)
		return fmt.Errorf("deployment returns an error: status code(%d) %s",
			response.StatusCode, string(openAPIErr.Body()))
	}
	var ret string
	if deployOptions.O != "" {
		ret, err = deployOptions.Output.Formatter(deployResp, err)
	} else {
		ret = printPlain(deployOptions, deployResp, statusOptions.DeploymentId, err)
	}

	fmt.Fprintln(cmd.OutOrStdout(), ret)

	return nil
}

func printPlain(deployOptions *cmd.RootOptions, deployResp deploy.DeploymentV2DeploymentStatusResponse, deploymentId string, err error) string {
	ret := ""
	if err != nil {
		logrus.Error(err)
		logrus.Fatalf("Error getting deployment status")
	}

	now := time.Now().Format(time.RFC3339)
	ret += fmt.Sprintf("[%v] application: %s, started: %s\n", now, deployResp.GetApplication(), deployResp.GetStartedAtIso8601())
	ret += fmt.Sprintf("[%v] status: ", now)
	switch status := deployResp.GetStatus(); status {
	case deploy.DEPLOYMENT_PAUSED:
		end := deployResp.Kubernetes.Canary.PauseInfo.GetEndTimeIso8601()
		reason := deployResp.Kubernetes.Canary.PauseInfo.GetReason()
		if reason == "" {
			reason = "unspecified"
		}
		ret += fmt.Sprintf("[%s] msg: Paused until %s for reason: %s. You may resume immediately in the cloud console or CLI\n", status, end, reason)
	case deploy.DEPLOYMENT_AWAITING_APPROVAL:
		ret += fmt.Sprintf("[%s] msg: Paused for Manual Judgment. You may resume immediately in the cloud console or CLI.\n", status)
	default:
		ret += string(status) + "\n"

	}
	url := cloudConsoleBaseUrl
	if strings.Contains(deployOptions.TokenIssuerUrl, "staging" ) {
		url = cloudConsoleStagingBaseUrl
	}
	url += "/deployments/" + deploymentId + "?environmentId=" + deployOptions.Environment
	ret += fmt.Sprintf("[%v] See the deployment status user interface: %s\n", now,url)
	return ret
}
