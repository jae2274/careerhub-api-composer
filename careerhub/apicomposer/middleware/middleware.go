package middleware

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jae2274/careerhub-api-composer/careerhub/apicomposer/jwtresolver"
)

func SetClaimsMW(jr *jwtresolver.JwtResolver) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("Authorization")
			if tokenString != "" {
				claims, err := jr.ParseToken(tokenString)
				if err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}

				if err := jr.Validate(claims); err != nil {
					w.WriteHeader(http.StatusUnauthorized)
				}

				ctx := context.WithValue(r.Context(), "claims", claims)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func GetClaims(ctx context.Context) (*jwtresolver.CustomClaims, bool) {
	claims, ok := ctx.Value("claims").(*jwtresolver.CustomClaims)
	return claims, ok
}
