package rollout

import (
	"context"
	"errors"
	"fmt"
	"github.com/armory/armory-cli/internal/deng/protobuff"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

const (
	ParamName = "name"
	ParamType = "type"
)

func Resume(ctx context.Context, cmd *cobra.Command, client protobuff.DeploymentServiceClient, args []string) error {
	return performRolloutOperation(ctx, cmd, args, client.Resume)
}

func Restart(ctx context.Context, cmd *cobra.Command, client protobuff.DeploymentServiceClient, args []string) error {
	return performRolloutOperation(ctx, cmd, args, client.Restart)
}

func Abort(ctx context.Context, cmd *cobra.Command, client protobuff.DeploymentServiceClient, args []string) error {
	return performRolloutOperation(ctx, cmd, args, client.Abort)
}

func performRolloutOperation(ctx context.Context, cmd *cobra.Command, args []string, call func(ctx context.Context, in *protobuff.RolloutRequest, opts ...grpc.CallOption) (*protobuff.RolloutResponse, error)) error {
	if len(args) == 0 {
		return errors.New("please provide deployment ID")
	}

	n, err := cmd.Flags().GetString(ParamName)
	if err != nil {
		return err
	}

	t, err := cmd.Flags().GetString(ParamType)
	if err != nil {
		return err
	}

	var req *protobuff.RolloutRequest
	depId := args[0]
	if n == "" {
		req = &protobuff.RolloutRequest{All: true, DeploymentId: depId}
	} else {
		req = &protobuff.RolloutRequest{Name: n, Type: t, DeploymentId: depId}
	}

	res, err := call(ctx, req)
	if err != nil {
		return err
	}
	fmt.Println(res.Message)
	return nil
}
