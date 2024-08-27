package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Log struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Action    string             `bson:"action"`
	Details   string             `bson:"details"`
	CreatedAt primitive.DateTime `bson:"created_at"`
}
