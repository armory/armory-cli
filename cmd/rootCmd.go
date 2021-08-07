package cmd

import (
	"context"
	"flag"
	"github.com/armory/armory-cli/internal/config"
	"github.com/armory/armory-cli/internal/deng"
	"github.com/armory/armory-cli/internal/helpers"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
)

const (
	deployCliName = "armory"
	ParamVerbose  = "verbose"
)

var verboseFlag bool

var RootCmd = &cobra.Command{
	Use:   deployCliName,
	Short: "Trigger, monitor, and diagnose your deployments",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(RootCmd.Execute())
}

func init() {
	RootCmd.PersistentFlags().BoolVarP(&verboseFlag, ParamVerbose, "v", false, "show more details")
	RootCmd.PersistentFlags().StringP(config.ParamContext, "C", "default", "context")

	// Hidden flags
	RootCmd.PersistentFlags().String(config.ParamEndpoint, "deploy.cloud.armory.io:443", "deploy engine endpoint")

	RootCmd.PersistentFlags().Bool(config.ParamInsecure, false, "do not verify server certificate")
	RootCmd.PersistentFlags().Bool(config.ParamPlaintext, false, "use a plaintext connection (warning insecure!)")
	RootCmd.PersistentFlags().Bool(config.ParamNoProxy, false, "skip system defined proxy (HTTP_PROXY, HTTPS_PROXY)")
	RootCmd.PersistentFlags().String(config.ParamCacert, "", "path to server certificate authority")
	RootCmd.PersistentFlags().String(config.ParamCert, "", "path to client certificate (mTLS)")
	RootCmd.PersistentFlags().String(config.ParamKey, "", "path to client certificate key (mTLS)")
	RootCmd.PersistentFlags().String(config.ParamKeyPassword, "", "password to the client certificate key (mTLS)")
	RootCmd.PersistentFlags().String(config.ParamServerName, "", "override server name")
	RootCmd.PersistentFlags().String(config.ParamToken, "", "authentication token")
	RootCmd.PersistentFlags().Bool(config.ParamAnonymously, false, "connect anonymously. This will likely fail in a non test environment.")

	RootCmd.PersistentPreRunE = configureLogging
	RootCmd.SilenceUsage = true
}

func configureLogging(cmd *cobra.Command, args []string) error {
	lvl := log.InfoLevel
	if verboseFlag {
		lvl = log.DebugLevel
	}
	log.SetLevel(lvl)
	log.SetFormatter(&log.TextFormatter{})
	_ = flag.Set("logtostderr", "true")
	return nil
}

type executor func(ctx context.Context, cmd *cobra.Command, client deng.DeploymentServiceClient, args []string) error

func ExecuteCancelable(cmd *cobra.Command, exe executor, args []string) error {
	ctx, cancel := context.WithCancel(context.TODO())

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, os.Interrupt)

	// Wait for signal to stop server
	go func() {
		<-signalCh
		log.Debug("signal received, stopping command...")
		cancel()
	}()

	dc, err := helpers.MakeDeploymentClient(ctx, cmd)
	if err != nil {
		return err
	}
	return exe(ctx, cmd, dc, args)
}