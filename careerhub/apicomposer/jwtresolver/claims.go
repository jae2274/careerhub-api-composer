package jwtresolver

import "github.com/golang-jwt/jwt/v5"

type CustomClaims struct {
	UserId string
	Roles  []string
	jwt.RegisteredClaims
}

func (c *CustomClaims) HasRole(role string) bool {
	for _, r := range c.Roles {
		if r == role {
			return true
		}
	}
	return false
}
