package api

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/mukeshmahato17/hrs/authutil"
	"github.com/mukeshmahato17/hrs/db"
	"github.com/mukeshmahato17/hrs/types"
	"go.mongodb.org/mongo-driver/mongo"
)

type HandleAuthStore struct {
	userStore db.UserStore
}

func NewHandleAuthStore(userStore db.UserStore) *HandleAuthStore {
	return &HandleAuthStore{
		userStore: userStore,
	}
}

type AuthParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	User  *types.User
	Token string
}

type genericResp struct {
	Type string `json:"type"`
	Msg  string `json:"msg"`
}

func invalidCredentials(c *fiber.Ctx) error {
	return c.Status(http.StatusBadRequest).JSON(genericResp{
		Type: "error",
		Msg:  "invalid credentials",
	})
}

func (h *HandleAuthStore) HandleAuthentication(c *fiber.Ctx) error {
	var params AuthParams
	if err := c.BodyParser(&params); err != nil {
		return err
	}

	user, err := h.userStore.GetUserByEmail(c.Context(), params.Email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return invalidCredentials(c)
		}
	}
	if !types.IsValidPassword(user.EncryptedPassword, params.Password) {
		return invalidCredentials(c)
	}
	resp := AuthResponse{
		User:  user,
		Token: CreateUserToken(user),
	}
	return c.JSON(resp)
}

func CreateUserToken(user *types.User) string {
	now := time.Now()
	expires := now.Add(time.Hour * 4).Unix()
	claims := jwt.MapClaims{
		"userID":  user.ID.Hex(),
		"email":   user.Email,
		"expires": expires,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := authutil.JWTSecret()

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("failed to sign token:", err)
		return ""
	}

	return tokenStr
}
