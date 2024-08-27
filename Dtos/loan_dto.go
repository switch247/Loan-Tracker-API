package Dtos

import "go.mongodb.org/mongo-driver/bson/primitive"

type GetLoan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Purpose   string             `bson:"purpose"`
	Status    string             `bson:"status" default:"pending"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

type UpdateLoan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Status    string             `bson:"status"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}

type CreateLoan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Purpose   string             `bson:"purpose"`
	Status    string             `bson:"status" default:"pending"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}
