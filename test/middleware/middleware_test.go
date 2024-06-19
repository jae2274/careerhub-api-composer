package middleware

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/jae2274/goutils/mw/httpmw"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	secretKey := "just a secret key"
	port := 33333
	rootURL := fmt.Sprintf("http://localhost:%d", port)

	getBody := func(res *http.Response) string {
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		return string(body)
	}

	t.Run("don't need jwt", func(t *testing.T) {
		t.Run("response 'Hello, World!' without jwt", func(t *testing.T) {
			res, err := initClient(t, secretKey, port).Get(rootURL + "/hello")
			require.NoError(t, err)

			defer res.Body.Close()
			require.Equal(t, "Hello, World!", getBody(res))
		})

		t.Run("response Hello ${userId} with jwt", func(t *testing.T) {
			req, err := http.NewRequest("GET", rootURL+"/hello", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

			res, err := initClient(t, secretKey, port).Do(req)
			require.NoError(t, err)

			defer res.Body.Close()
			require.Equal(t, "Hello, Jyo Liar!", getBody(res))
		})
	})

	t.Run("need login", func(t *testing.T) {
		t.Run("response 401 without jwt", func(t *testing.T) {
			res, err := initClient(t, secretKey, port).Get(rootURL + "/my")
			require.NoError(t, err)
			defer res.Body.Close()

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			require.Equal(t, "", getBody(res))
		})

		t.Run("response 200 with jwt", func(t *testing.T) {
			req, err := http.NewRequest("GET", rootURL+"/my", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

			res, err := initClient(t, secretKey, port).Do(req)
			require.NoError(t, err)

			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, "Hello, Jyo Liar!", getBody(res))
		})
	})

	t.Run("need role: admin", func(t *testing.T) {
		t.Run("response 401 without jwt", func(t *testing.T) {
			res, err := initClient(t, secretKey, port).Get(rootURL + "/admin")
			require.NoError(t, err)
			defer res.Body.Close()

			require.Equal(t, http.StatusUnauthorized, res.StatusCode)
			require.Equal(t, "", getBody(res))
		})

		t.Run("response 403 with jwt but without role", func(t *testing.T) {
			req, err := http.NewRequest("GET", rootURL+"/admin", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

			res, err := initClient(t, secretKey, port).Do(req)
			require.NoError(t, err)

			defer res.Body.Close()
			require.Equal(t, http.StatusForbidden, res.StatusCode)
			require.Equal(t, "", getBody(res))
		})

		t.Run("response 200 with jwt and role", func(t *testing.T) {
			req, err := http.NewRequest("GET", rootURL+"/admin", nil)
			require.NoError(t, err)
			req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{"admin"}))

			res, err := initClient(t, secretKey, port).Do(req)
			require.NoError(t, err)

			defer res.Body.Close()
			require.Equal(t, http.StatusOK, res.StatusCode)
			require.Equal(t, "Hello, Jyo Liar!", getBody(res))
		})
	})

	t.Run("need role for specify url path: manager", func(t *testing.T) {
		t.Run("all response 401 without jwt", func(t *testing.T) {
			testClient := initClient(t, secretKey, port)

			for _, path := range []string{"/manager", "/manager/a", "/manager/a/b", "/manager/a/b/c"} {
				res, err := testClient.Get(rootURL + path)
				require.NoError(t, err)
				defer res.Body.Close()

				require.Equal(t, http.StatusUnauthorized, res.StatusCode)
				require.Equal(t, "", getBody(res))
			}
		})

		t.Run("response 403 with jwt but without role", func(t *testing.T) {
			testClient := initClient(t, secretKey, port)

			for _, path := range []string{"/manager", "/manager/a", "/manager/a/b", "/manager/a/b/c"} {
				req, err := http.NewRequest("GET", rootURL+path, nil)
				require.NoError(t, err)
				req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

				res, err := testClient.Do(req)
				require.NoError(t, err)

				defer res.Body.Close()
				require.Equal(t, http.StatusForbidden, res.StatusCode)
				require.Equal(t, "", getBody(res))
			}
		})

		t.Run("response 200 with jwt and role", func(t *testing.T) {
			testClient := initClient(t, secretKey, port)

			for _, path := range []string{"/manager", "/manager/a", "/manager/a/b", "/manager/a/b/c"} {
				req, err := http.NewRequest("GET", rootURL+path, nil)
				require.NoError(t, err)
				req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{"manager"}))

				res, err := testClient.Do(req)
				require.NoError(t, err)

				defer res.Body.Close()
				require.Equal(t, http.StatusOK, res.StatusCode)
				require.Equal(t, "Hello, Jyo Liar!", getBody(res))
			}
		})
	})

}

func initRouter() *mux.Router {
	jr := jwtresolver.NewJwtResolver("just a secret key")
	commonHandleFunc := func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if ok {
			w.Write([]byte("Hello, " + claims.UserId + "!"))
			return
		}
		w.Write([]byte("Hello, World!"))
	}

	router := mux.NewRouter()
	router.Use(httpmw.SetTraceIdMW(), middleware.SetClaimsMW(jr))
	router.HandleFunc("/hello", commonHandleFunc)

	myRouter := router.PathPrefix("/my").Subrouter()
	myRouter.Use(middleware.CheckJustLoggedIn)
	myRouter.HandleFunc("", commonHandleFunc)

	adminRouter := router.PathPrefix("/admin").Subrouter()
	adminRouter.Use(middleware.CheckHasRole("admin"))
	adminRouter.HandleFunc("", commonHandleFunc)

	managerRouter := router.PathPrefix("/manager").Subrouter()
	managerRouter.Use(middleware.CheckHasRole("manager"))
	managerRouter.HandleFunc("", commonHandleFunc)
	managerRouter.HandleFunc("/a", commonHandleFunc)
	managerRouter.HandleFunc("/a/b", commonHandleFunc)
	managerRouter.HandleFunc("/a/b/c", commonHandleFunc)

	return router
}

func initClient(t *testing.T, secretKey string, port int) *http.Client {
	router := initRouter()

	go func() {
		http.ListenAndServe(fmt.Sprintf(":%d", port), router)
	}()

	return &http.Client{}
}

func createAccessToken(t *testing.T, secretKey, userId string, authorities []string) string {
	now := time.Now()

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&jwtresolver.CustomClaims{
			UserId:      userId,
			Authorities: authorities,
			RegisteredClaims: jwt.RegisteredClaims{
				Issuer:    "careerhub.jyo-liar.com",                    //TODO: 임의 설정
				Audience:  []string{"careerhub.jyo-liar.com"},          //TODO: 임의 설정
				ExpiresAt: jwt.NewNumericDate(now.Add(30 * time.Hour)), //TODO: 임의 설정
				IssuedAt:  jwt.NewNumericDate(now),
				NotBefore: jwt.NewNumericDate(now),
			},
		},
	).SignedString([]byte(secretKey))

	require.NoError(t, err)

	return accessToken
}
