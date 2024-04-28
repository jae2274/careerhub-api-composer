package jwtresolver

import (
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/terr"
)

type CustomClaims struct {
	UserId string
	Roles  []string
	jwt.RegisteredClaims
}

type JwtResolver struct {
	secretKey []byte
	validator *jwt.Validator
}

type TokenInfo struct {
	GrantType    string `json:"grantType"`
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

func NewJwtResolver(secretKey string) *JwtResolver {
	return &JwtResolver{
		secretKey: []byte(secretKey),
		validator: jwt.NewValidator(),
	}
}

func (j *JwtResolver) ParseToken(tokenString string) (*CustomClaims, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	jwt, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := jwt.Claims.(*CustomClaims); ok {
		return claims, nil
	} else {
		return nil, terr.New("invalid token. claims is not CustomClaims type")
	}
}

func (j *JwtResolver) Validate(claims *CustomClaims) error {
	return j.validator.Validate(claims)
}

func (j *JwtResolver) CheckHasRole(role string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
			claims, err := j.ParseToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if err := j.Validate(claims); err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			if slices.Contains(claims.Roles, role) {
				next.ServeHTTP(w, r)
			} else {
				w.WriteHeader(http.StatusForbidden)
			}
		})
	}
}
