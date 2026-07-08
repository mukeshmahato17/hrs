package api

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/types"
)

func getAuthUser(c *fiber.Ctx) (*types.User, error) {
	user, ok := c.Context().UserValue("user").(*types.User)
	if !ok {
		return nil, fmt.Errorf("unauthorized user")
	}
	return user, nil
}
