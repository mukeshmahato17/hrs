package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/types"
)

func HandleUsers(c *fiber.Ctx) error {
	return c.JSON("users")
}

func HandleUser(c *fiber.Ctx) error {
	user := types.User{
		FirstName: "Foo",
		LastName:  "Bar",
	}
	return c.JSON(user)
}
