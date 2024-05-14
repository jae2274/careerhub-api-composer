package vars

import (
	"fmt"
	"os"
	"strconv"
)

type Vars struct {
	PostingGrpcEndpoint      string
	UserinfoGrpcEndpoint     string
	ReviewGrpcEndpoint       string
	RootPath                 string
	SecretKey                string
	AccessControlAllowOrigin *string
	ApiPort                  int
}

type ErrNotExistedVar struct {
	VarName string
}

func NotExistedVar(varName string) *ErrNotExistedVar {
	return &ErrNotExistedVar{VarName: varName}
}

func (e *ErrNotExistedVar) Error() string {
	return fmt.Sprintf("%s is not existed", e.VarName)
}

func Variables() (*Vars, error) {
	postingGrpcEndpoint, err := getFromEnv("POSTING_GRPC_ENDPOINT")
	if err != nil {
		return nil, err
	}

	userinfoGrpcEndpoint, err := getFromEnv("USERINFO_GRPC_ENDPOINT")
	if err != nil {
		return nil, err
	}

	reviewGrpcEndpoint, err := getFromEnv("REVIEW_GRPC_ENDPOINT")
	if err != nil {
		return nil, err
	}

	secretKey, err := getFromEnv("SECRET_KEY")
	if err != nil {
		return nil, err
	}

	apiPort, err := getFromEnv("API_PORT")
	if err != nil {
		return nil, err
	}

	apiPortInt, err := strconv.ParseInt(apiPort, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("API_PORT is not integer.\tAPI_PORT: %s", apiPort)
	}

	accesControlAllowOrigin := getFromEnvPtr("ACCESS_CONTROL_ALLOW_ORIGIN")

	return &Vars{
		PostingGrpcEndpoint:      postingGrpcEndpoint,
		UserinfoGrpcEndpoint:     userinfoGrpcEndpoint,
		ReviewGrpcEndpoint:       reviewGrpcEndpoint,
		RootPath:                 os.Getenv("ROOT_PATH"),
		SecretKey:                secretKey,
		ApiPort:                  int(apiPortInt),
		AccessControlAllowOrigin: accesControlAllowOrigin,
	}, nil
}

func getFromEnv(envVar string) (string, error) {
	ev := os.Getenv(envVar)

	if ev == "" {
		return "", fmt.Errorf("%s is not existed", envVar)
	}

	return ev, nil
}

func getFromEnvPtr(envVar string) *string {
	ev := os.Getenv(envVar)

	if ev == "" {
		return nil
	}

	return &ev
}
