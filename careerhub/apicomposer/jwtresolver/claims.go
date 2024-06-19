package jwtresolver

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId      string `validate:"nonzero"`
	Authorities []string
	jwt.RegisteredClaims
}

func (c *CustomClaims) HasAuthority(authority string) bool {
	for _, r := range c.Authorities {
		if r == authority {
			return true
		}
	}
	return false
}
