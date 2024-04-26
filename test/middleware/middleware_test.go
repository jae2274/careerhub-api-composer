package middleware

import (
	"io"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/middleware"
	"github.com/stretchr/testify/require"
)

func TestMiddleware(t *testing.T) {
	secretKey := "just a secret key"
	rootURL := "http://localhost:8080"
	t.Run("response 'Hello, World!' without jwt", func(t *testing.T) {
		res, err := initClient(t, secretKey).Get(rootURL + "/hello")
		require.NoError(t, err)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "Hello, World!", string(body))
	})

	t.Run("response Hello ${userId} with jwt", func(t *testing.T) {
		req, err := http.NewRequest("GET", rootURL+"/hello", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

		res, err := initClient(t, secretKey).Do(req)
		require.NoError(t, err)

		defer res.Body.Close()
		body, err := io.ReadAll(res.Body)
		require.NoError(t, err)
		require.Equal(t, "Hello, Jyo Liar!", string(body))
	})

	t.Run("response 401 if call admin api and without jwt", func(t *testing.T) {
		res, err := initClient(t, secretKey).Get(rootURL + "/admin")
		require.NoError(t, err)
		defer res.Body.Close()

		require.Equal(t, http.StatusUnauthorized, res.StatusCode)
	})

	t.Run("response 403 if call admin api and with jwt but without role", func(t *testing.T) {
		req, err := http.NewRequest("GET", rootURL+"/admin", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{}))

		res, err := initClient(t, secretKey).Do(req)
		require.NoError(t, err)

		defer res.Body.Close()
		require.Equal(t, http.StatusForbidden, res.StatusCode)
	})

	t.Run("response 200 if call admin api and with jwt and role", func(t *testing.T) {
		req, err := http.NewRequest("GET", rootURL+"/admin", nil)
		require.NoError(t, err)
		req.Header.Set("Authorization", createAccessToken(t, secretKey, "Jyo Liar", []string{"admin"}))

		res, err := initClient(t, secretKey).Do(req)
		require.NoError(t, err)

		defer res.Body.Close()
		require.Equal(t, http.StatusOK, res.StatusCode)
	})
}

func initRouter() *mux.Router {
	router := mux.NewRouter()
	jr := jwtresolver.NewJwtResolver("just a secret key")
	router.Use(middleware.SetClaimsMW(jr))

	router.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if ok {
			w.Write([]byte("Hello, " + claims.UserId + "!"))
			return
		}
		w.Write([]byte("Hello, World!"))
	})

	router.HandleFunc("/admin", func(w http.ResponseWriter, r *http.Request) {
		claims, ok := middleware.GetClaims(r.Context())
		if !ok {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		if !slices.Contains(claims.Roles, "admin") {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	return router
}

func initClient(t *testing.T, secretKey string) *http.Client {
	router := initRouter()

	go func() {
		http.ListenAndServe(":8080", router)
	}()

	return &http.Client{}
}

func createAccessToken(t *testing.T, secretKey, userId string, roles []string) string {
	now := time.Now()

	accessToken, err := jwt.NewWithClaims(jwt.SigningMethodHS256,
		&jwtresolver.CustomClaims{
			UserId: userId,
			Roles:  roles,
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
