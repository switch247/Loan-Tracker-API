package Dtos

import (
	"reflect"
	"time"

	"github.com/go-playground/validator"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RegisterUserDto struct {
	Email          string    `json:"email" validate:"required,email"`
	Password       string    `json:"password" validate:"required"`
	UserName       string    `json:"username" `
	Role           string    `json:"-",omitempty default:"user"`
	ProfilePicture string    `json:"profile_picture"`
	Bio            string    `json:"bio"`
	EmailVerified  bool      `bson:"email_verified" default:"false"`
	Name           string    `json:name`
	CreatedAt      time.Time `json:"createdat"`
	UpdatedAt      time.Time `json:"updatedat"`
}

type LoginUserDto struct {
	Email    string `json:"email" validate:"required,email"`
	UserName string `json:"username",omitempty`
	Password string `json:"password" validate:"required"`
}

func CustomValidator(fl validator.FieldLevel) bool {
	// Get the struct field value
	v := fl.Field()

	// Check if both Email and UserName are not empty
	switch v.Type().Kind() {
	case reflect.Struct:
		email := v.FieldByName("Email").String()
		username := v.FieldByName("UserName").String()
		return email != "" || username != ""
	default:
		return false
	}
}

// this could have been handled in a better way but i was too lazy to do it
type OmitedUser struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	UserName       string             `bson:"username"`
	Email          string             `json:"email" validate:"required"`
	Password       string             `json:"-"`
	Role           string             `json:"role"`
	ProfilePicture string             `json:"profile_picture"`
	Bio            string             `json:"bio"`
	EmailVerified  bool               `bson:"email_verified"`
	Name           string             `json:name`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}

type UpdateUser struct {
	Name     string `json:"name"`
	UserName string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Bio      string `json:"bio"`
}
