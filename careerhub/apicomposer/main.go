package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/controller"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting"
	postingGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/posting/restapi_grpc"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo"
	userinfoGrpc "github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/userinfo/restapi_grpc"
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

	conn, err := grpc.NewClient(envVars.PostingGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(mainCtx, err)
	postingClient := postingGrpc.NewRestApiGrpcClient(conn)

	conn, err = grpc.NewClient(envVars.UserinfoGrpcEndpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	checkErr(mainCtx, err)
	matchJobClient := userinfoGrpc.NewMatchJobGrpcClient(conn)

	jr := jwtresolver.NewJwtResolver(envVars.SecretKey)

	//rootRouter 설정
	rootRouter := mux.NewRouter().PathPrefix(envVars.RootPath).Subrouter()

	rootRouter.Use(httpmw.SetTraceIdMW(), middleware.SetClaimsMW(jr))

	controller.NewJobPostingController(
		posting.NewPostingService(postingClient),
	).RegisterRoutes(rootRouter)

	//userinfoRouter 설정
	userinfoRouter := rootRouter.PathPrefix("/my").Subrouter()
	userinfoRouter.Use(middleware.CheckJustLoggedIn)

	controller.NewUserinfoController(
		userinfo.NewMatchJobService(matchJobClient),
	).RegisterRoutes(userinfoRouter)

	var allowOrigins []string
	if envVars.AccessControlAllowOrigin != nil {
		allowOrigins = append(allowOrigins, *envVars.AccessControlAllowOrigin)
	}
	originsOk := handlers.AllowedOrigins(allowOrigins)
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	llog.Msg("Start composed api server").Level(llog.INFO).Data("port", envVars.ApiPort).Data("rootPath", envVars.RootPath).Log(mainCtx)

	err = http.ListenAndServe(fmt.Sprintf(":%d", envVars.ApiPort), handlers.CORS(originsOk, headersOk, methodsOk)(rootRouter))
	checkErr(mainCtx, err)
}

func checkErr(ctx context.Context, err error) {
	if err != nil {
		llog.LogErr(ctx, err)
		os.Exit(1)
	}
}
