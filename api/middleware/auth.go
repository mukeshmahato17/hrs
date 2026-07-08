package middleware

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mukeshmahato17/hrs/authutil"
	"github.com/mukeshmahato17/hrs/db"
)

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if token == "" {
			authorization := c.Get("Authorization")
			if strings.HasPrefix(authorization, "Bearer ") {
				token = strings.TrimSpace(strings.TrimPrefix(authorization, "Bearer "))
			}
		}
		if token == "" {
			return fmt.Errorf("unauthorized")
		}
		claims, err := validateToken(token)
		if err != nil {
			return err
		}

		expiresValue, ok := claims["expires"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		expiresFloat, ok := expiresValue.(float64)
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		if time.Now().Unix() > int64(expiresFloat) {
			return fmt.Errorf("token expired")
		}

		userIDValue, ok := claims["userID"]
		if !ok {
			return fmt.Errorf("unauthorized")
		}
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return fmt.Errorf("unauthorized")
		}

		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return fmt.Errorf("unauthorized")
		}

		c.Context().SetUserValue("user", user)
		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("imvalid signing method", token.Header["alg"])
			return nil, fmt.Errorf("unauthorized")
		}
		secret := authutil.JWTSecret()
		return []byte(secret), nil
	})
	if err != nil {
		fmt.Println("failed to parse token:", err)
		return nil, fmt.Errorf("unauthorized")
	}
	if !token.Valid {
		fmt.Println("invalid token")
		return nil, fmt.Errorf("unauthorize")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("unauthorize")
	}
	return claims, nil
}
