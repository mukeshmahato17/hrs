package authutil

import "os"

const DefaultJWTSecret = "dev-jwt-secret"

func JWTSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret != "" {
		return secret
	}
	return DefaultJWTSecret
}
