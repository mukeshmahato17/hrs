package types

import (
	"fmt"
	"regexp"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCode      = 12
	minFirstNameLen = 2
	minLastNameLen  = 2
	minPasswordLen  = 6
)

var EmailRegexp = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type CreateUserParams struct {
	FirstName string `json: "firstName"`
	LastName  string `json: "lastName"`
	Email     string `json: "email"`
	Password  string `json: "password"`
}

func (params CreateUserParams) Validate() []string {
	var errors []string
	if len(params.FirstName) < minFirstNameLen {
		errors = append(errors, fmt.Sprintf("first name length should be at least %d characters", minFirstNameLen))
	}
	if len(params.LastName) < minLastNameLen {
		errors = append(errors, fmt.Sprintf("last name length should be at least %d characters", minLastNameLen))
	}
	if len(params.Password) < minPasswordLen {
		errors = append(errors, fmt.Sprintf("password length should be at least %d characters", minPasswordLen))
	}

	if !IsValidEmail(params.Email) {
		errors = append(errors, fmt.Sprintf("invalid email"))
	}
	return errors
}

// IsValidEmail checks if a string matches the email regex format
func IsValidEmail(email string) bool {
	return EmailRegexp.MatchString(email)
}

type User struct {
	ID                primitive.ObjectID `bson:"_id,omitempty", json: "id,omitempty"`
	FirstName         string             `bson: "firstName", json: "firstName"`
	LastName          string             `bson: "lastName", json: "lastName"`
	Email             string             `bson: "email", json: "email"`
	EncryptedPassword string             `bson: "EncryptedPassword", json: "-"`
}

func NewUserFromParams(params CreateUserParams) (*User, error) {
	encpw, err := bcrypt.GenerateFromPassword([]byte(params.Password), bcryptCode)
	if err != nil {
		return nil, err
	}

	return &User{
		FirstName:         params.FirstName,
		LastName:          params.LastName,
		Email:             params.Email,
		EncryptedPassword: string(encpw),
	}, nil
}
