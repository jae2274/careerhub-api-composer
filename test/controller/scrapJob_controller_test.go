package controller

// import (
// 	"net/http"
// 	"testing"

// 	"github.com/jae2274/careerhub-api-composer/test/tinit"
// )

// func TestScrapJobController(t *testing.T) {
// 	cancelFunc := tinit.RunTestApp(t)
// 	defer cancelFunc()

// 	testHost := tinit.InitTestHost(t)

// 	requests := []*http.Request{
// 		{
// 			Method: "GET",
// 			URL:    testHost.GetUrl(t, "/scrap-job"),
// 		},
// 		{
// 			Method: "POST",
// 			URL:    testHost.GetUrl(t, "/scrap-job"),
// 		},
// 		{
// 			Method: "DELETE",
// 			URL:    testHost.GetUrl(t, "/scrap-job"),
// 		},
// 		{
// 			Method: "GET",
// 			URL:    testHost.GetUrl(t, "/scrap-job/untagged"),
// 		},
// 		{
// 			Method: "GET",
// 			URL:    testHost.GetUrl(t, "/scrap-job/tags"),
// 		},
// 		{
// 			Method: "POST",
// 			URL:    testHost.GetUrl(t, "/scrap-job/tags"),
// 		},
// 		{
// 			Method: "DELETE",
// 			URL:    testHost.GetUrl(t, "/scrap-job/tags"),
// 		},
// 	}

// }
