package vars

import (
	"fmt"
	"os"
	"strconv"
)

type Vars struct {
	PostingGrpcEndpoint string
	RootPath            string
	ApiPort             int
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

	apiPort, err := getFromEnv("API_PORT")
	if err != nil {
		return nil, err
	}

	apiPortInt, err := strconv.ParseInt(apiPort, 10, 32)
	if err != nil {
		return nil, fmt.Errorf("API_PORT is not integer.\tAPI_PORT: %s", apiPort)
	}

	return &Vars{
		PostingGrpcEndpoint: postingGrpcEndpoint,
		RootPath:            os.Getenv("ROOT_PATH"),
		ApiPort:             int(apiPortInt),
	}, nil
}

func getFromEnv(envVar string) (string, error) {
	ev := os.Getenv(envVar)

	if ev == "" {
		return "", fmt.Errorf("%s is not existed", envVar)
	}

	return ev, nil
}
