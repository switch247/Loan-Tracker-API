package Domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type Loan struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	UserID    primitive.ObjectID `bson:"user_id"`
	Amount    float64            `bson:"amount"`
	Purpose   string             `bson:"purpose"`
	Status    string             `bson:"status"`
	CreatedAt primitive.DateTime `bson:"created_at"`
	UpdatedAt primitive.DateTime `bson:"updated_at"`
}
