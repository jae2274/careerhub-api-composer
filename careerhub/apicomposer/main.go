package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/vars"
	"github.com/jae2274/goutils/llog"
	"github.com/jae2274/goutils/mw"
	"github.com/jae2274/goutils/mw/httpmw"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	app     = "api-composer"
	service = "careerhub"

	ctxKeyTraceID = string(mw.CtxKeyTraceID)

	// needRole = "ROLE_CAREERHUB_USER"
)

func initLogger(ctx context.Context) error {
	llog.SetMetadata("service", service)
	llog.SetMetadata("app", app)
	llog.SetDefaultContextData(ctxKeyTraceID)

	hostname, err := os.Hostname()
	if err != nil {
		return err
	}

	llog.SetMetadata("hostname", hostname)

	return nil
}

func main() {
	mainCtx := context.Background()

	err := initLogger(mainCtx)
	checkErr(mainCtx, err)
	llog.Info(mainCtx, "Start Application")

	envVars, err := vars.Variables()
	checkErr(mainCtx, err)

	conn, err := grpc.Dial(envVars.PostingGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(mainCtx, err)

	postingClient := restapi_grpc.NewRestApiGrpcClient(conn)

	router := mux.NewRouter()
	// jwtResolver := jwtresolver.NewJwtResolver(envVars.SecretKey)
	// authMiddleware := jwtResolver.CheckHasRole(needRole)
	router.Use(
		// authMiddleware,
		httpmw.SetTraceIdMW()) //TODO: 불필요한 파라미터가 잘못 포함되어 있어 이후 라이브러리 수정 필요
	ctrler := NewController(
		posting.NewPostingService(postingClient),
		router,
	)

	ctrler.RegisterRoutes(envVars.RootPath)

	var allowOrigins []string
	if envVars.AccessControlAllowOrigin != nil {
		allowOrigins = append(allowOrigins, *envVars.AccessControlAllowOrigin)
	}
	originsOk := handlers.AllowedOrigins(allowOrigins)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	llog.Msg("Start composed api server").Level(llog.INFO).Data("port", envVars.ApiPort).Data("rootPath", envVars.RootPath).Log(mainCtx)

	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), handlers.CORS(originsOk, headersOk, methodsOk)(router))
	checkErr(mainCtx, err)
}

func checkErr(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		os.Exit(1)
	}
}
