package deployCmd

import (
	"github.com/armory/armory-cli/internal/deploy"
	"github.com/armory/armory-cli/internal/helpers"
	"github.com/spf13/cobra"
)

// For now we only map to a single artifact
// which could become multiple (e.g. kustomization, helm)
// Later we could read from a deployment file definition or even
// keep adding to a server-side deployment before kicking it off.
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Initiate a deployment to a target environment",
	Long: `Initiate and optionally monitor the deployment of artifacts to a target environment.
Environments are made available by agents - see https://github.com/armory-io/armory-agents`,
	RunE: func(c *cobra.Command, args []string) error {
		return helpers.ExecuteCancelable(c, deploy.Execute, args)
	},
}

func init() {
	startCmd.Flags().BoolP(deploy.ParameterWait, "w", false, "wait for deployment success or failure")
	startCmd.Flags().String(deploy.ParameterEnvironmentType, "kubernetes", "deployment account type")
	startCmd.Flags().StringP(deploy.ParameterEnvironmentName, "a", "", "deployment account name")
	startCmd.Flags().StringP(deploy.ParameterEnvironmentNamespace, "n", "", "(Kubernetes only) namespace to deploy to. Defaults to manifest namespace.")
	startCmd.Flags().String(deploy.ParameterViaAccount, "", "use specified agent to retrieve artifact")
	startCmd.Flags().String(deploy.ParameterViaProvider, "", "use agent of specified provider to retrieve artifact")
	startCmd.Flags().StringSlice(deploy.ParameterVersion, nil, "(Kubernetes only) specific container versions to override")
	startCmd.Flags().BoolP(deploy.ParameterKustomize, "k", false, "(Kubernetes only) parameter is a Kustomization")
	startCmd.Flags().BoolP(deploy.ParameterLocal, "l", false, "resolve artifacts locally")
	startCmd.Flags().String(deploy.ParameterApplication, "", "application this deployment is part of")
	startCmd.Flags().StringP(deploy.ParameterStrategy, "s", "update", "Strategy one of update,bluegreen,canary")
	startCmd.Flags().StringArray(deploy.ParameterStrategySteps, nil, "wait(duration), pause, ratio(valueOrPct), traffic(percent)")
	_ = startCmd.MarkFlagRequired(deploy.ParameterEnvironmentName)
}