package jwtresolver

import (
	"errors"
	"net/http"
	"slices"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gorilla/mux"
	"github.com/jae2274/goutils/terr"
	"gopkg.in/validator.v2"
)

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

func (j *JwtResolver) ParseToken(tokenString string) (*CustomClaims, bool, error) {

	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	jwtToken, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return j.secretKey, nil
	})

	if jwtToken.Valid {
		if claims, ok := jwtToken.Claims.(*CustomClaims); ok {
			if err := validator.Validate(claims); err != nil {
				return claims, false, err
			} else {
				return claims, true, nil
			}
		} else {
			return &CustomClaims{}, false, terr.New("invalid token. claims is not CustomClaims type")
		}
	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		return &CustomClaims{}, false, nil
	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		return &CustomClaims{}, false, terr.New("invalid token. token is malformed")
	} else {
		return &CustomClaims{}, false, terr.Wrap(err)
	}
}

func (j *JwtResolver) Validate(claims *CustomClaims) error {
	if err := validator.Validate(claims); err != nil {
		return err
	}

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
			claims, isValid, err := j.ParseToken(tokenString)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			if !isValid {
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
