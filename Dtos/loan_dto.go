package Dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLoan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Purpose   string             `bson:"purpose"`
	Status    string             `bson:"status" validate:"required,oneof=pending approved rejected" default:"pending"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

type UpdateLoan struct {
	Status    string             `bson:"status" validate:"required,oneof=pending approved rejected" default:"pending"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

// : pending | approved | rejected
type CreateLoan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Purpose   string             `bson:"purpose"`
	Status    string             `bson:"status" validate:"oneof=pending approved rejected" default:"pending"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}
