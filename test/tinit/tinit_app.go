package tinit

import (
	"context"
	"fmt"
	"testing"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/app"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/vars"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func RunTestApp(t *testing.T) context.CancelFunc {
	ctx, cancelFunc := context.WithCancel(context.Background())
	go app.Run(ctx)

	return cancelFunc
}
func InitEnvVars(t *testing.T) *vars.Vars {
	envVars, err := vars.Variables()
	require.NoError(t, err)

	return envVars
}

func InitGrpcClient(t *testing.T, port int) *grpc.ClientConn {
	target := fmt.Sprintf("localhost:%d", port)
	client, err := grpc.NewClient(target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	return client
}
