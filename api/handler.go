package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/types"
)

type HandleUserStore struct {
	userStore db.UserStore
}

func NewHandleUserStore(userStore db.UserStore) *HandleUserStore {
	return &HandleUserStore{
		userStore: userStore,
	}
}

func (h *HandleUserStore) HandleUserPost(c *fiber.Ctx) error {
	var params types.CreateUserParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if errors := params.Validate(); len(errors) > 0 {
		return c.JSON(errors)
	}
	user, err := types.NewUserFromParams(params)
	if err != nil {
		return err
	}
	insertedUser, err := h.userStore.InsertUser(c.Context(), *user)
	if err != nil {
		return err
	}
	return c.JSON(insertedUser)
}

func (h *HandleUserStore) HandleUsers(c *fiber.Ctx) error {
	users, err := h.userStore.GetUsers(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "internal server error fetching users",
		})
	}
	return c.JSON(users)
}

func (h *HandleUserStore) HandleUser(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := h.userStore.GetUserByID(c.Context(), id)
	if err != nil {
		return err
	}
	return c.JSON(user)
}
