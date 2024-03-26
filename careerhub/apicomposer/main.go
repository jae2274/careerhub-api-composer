package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/vars"
	"github.com/jae2274/goutils/llog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	needRole = "ROLE_CAREERHUB_USER"
)

func main() {
	mainCtx := context.Background()
	envVars, err := vars.Variables()
	checkErr(mainCtx, err)

	conn, err := grpc.Dial(envVars.PostingGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(mainCtx, err)

	postingClient := restapi_grpc.NewRestApiGrpcClient(conn)

	router := mux.NewRouter()
	jwtResolver := jwtresolver.NewJwtResolver(envVars.SecretKey)
	authMiddleware := jwtResolver.CheckHasRole(needRole)
	router.Use(authMiddleware)
	ctrler := NewController(postingClient, router)

	ctrler.RegisterRoutes(envVars.RootPath)

	llog.Msg("Composed api server is running").Level(llog.INFO).Data("apiPort", envVars.ApiPort).Data("rootPath", envVars.RootPath).Log(mainCtx)

	var allowOrigins []string
	if envVars.AccessControlAllowOrigin != nil {
		allowOrigins = append(allowOrigins, *envVars.AccessControlAllowOrigin)
	}
	originsOk := handlers.AllowedOrigins(allowOrigins)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), handlers.CORS(originsOk, headersOk, methodsOk)(router))
	checkErr(mainCtx, err)
}

func checkErr(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		panic(err)
	}
}
