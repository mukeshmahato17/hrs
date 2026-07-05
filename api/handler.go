package api

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandleUserStore struct {
	userStore db.UserStore
}

func NewHandleUserStore(userStore db.UserStore) *HandleUserStore {
	return &HandleUserStore{
		userStore: userStore,
	}
}

func (h *HandleUserStore) HandlePutUser(c *fiber.Ctx) error {
	var (
		// values = bson.M{}
		params = types.UpdateUserParams{}
		userID = c.Params("id")
	)
	oid, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return err
	}
	if err := c.BodyParser(&params); err != nil {
		return err
	}
	if err := h.userStore.UpdateUser(c.Context(), bson.M{"_id": oid}, params); err != nil {
		return err
	}
	return c.JSON(map[string]string{"updated": userID})
}

func (h *HandleUserStore) HandleDeleteUser(c *fiber.Ctx) error {
	userID := c.Params("id")
	if err := h.userStore.DeleteUser(c.Context(), userID); err != nil {
		return err
	}
	return c.JSON(map[string]string{"delete": userID})

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
		if errors.Is(err, mongo.ErrNoDocuments) {
			return c.JSON(map[string]string{"error": "not found"})
		}
	}
	return c.JSON(user)
}
