package Domain

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	ID             primitive.ObjectID `bson:"_id,omitempty" bson:"_id,omitempty"`
	Name           string             `json:"name"`
	UserName       string             `json:"username"`
	Email          string             `json:"email" validate:"required"`
	Password       string             `json:"password,omitempty" validate:"required"`
	Role           string             `json:"role"`
	ProfilePicture string             `json:"profile_picture"`
	Bio            string             `json:"bio"`
	EmailVerified  bool               `bson:"email_verified"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedAt      time.Time          `json:"updated_at"`
}
