package middleware

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/types"
)

func AdminAuth(c *fiber.Ctx) error {
	user, ok := c.Context().Value("user").(*types.User)
	if !ok {
		return fmt.Errorf("unauthorized")
	}
	if !user.IsAdmin {
		return fmt.Errorf("unauthorized")
	}
	return c.Next()
}
