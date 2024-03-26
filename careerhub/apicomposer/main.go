package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/vars"
	"github.com/jae2274/goutils/llog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	mainCtx := context.Background()
	envVars, err := vars.Variables()
	checkErr(mainCtx, err)

	conn, err := grpc.Dial(envVars.PostingGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(mainCtx, err)

	postingClient := restapi_grpc.NewRestApiGrpcClient(conn)

	router := mux.NewRouter()
	ctrler := NewController(postingClient, router)
	ctrler.RegisterRoutes(envVars.RootPath)

	llog.Msg("Composed api server is running").Level(llog.INFO).Data("apiPort", envVars.ApiPort).Data("rootPath", envVars.RootPath).Log(mainCtx)
	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), router)
	checkErr(mainCtx, err)
}

func checkErr(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
