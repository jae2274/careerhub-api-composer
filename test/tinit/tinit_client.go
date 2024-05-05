package tinit

import (
	"fmt"
	"net/url"
	"testing"

	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/vars"
	"github.com/stretchr/testify/require"
)

type TestHost struct {
	Protocol string
	Host     string
	Port     int
}

func (th *TestHost) GetUrl(t *testing.T, path string) *url.URL {
	u, err := url.Parse(fmt.Sprintf("%s://%s:%d%s", th.Protocol, th.Host, th.Port, path))
	require.NoError(t, err)

	return u
}

func InitTestHost(t *testing.T) *TestHost {
	envVars, err := vars.Variables()
	require.NoError(t, err)

	return &TestHost{
		Protocol: "http",
		Host:     "localhost",
		Port:     envVars.ApiPort,
	}
}
