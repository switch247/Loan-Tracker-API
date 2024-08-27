package routers

import (
	custommongo "Loan-Tracker-API/CustomMongo"
	"Loan-Tracker-API/Domain"
	"context"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Router *gin.Engine
var BlogCollections Domain.BlogCollections

func Setuprouter(client *mongo.Client) *gin.Engine {
	// Initialize the database
	DataBase := client.Database("loan-tracker")

	//initialize the user collections
	usercol := DataBase.Collection("Users")
	refreshtokencol := DataBase.Collection("RefreshTokens")

	// Initialize the custom user collections
	customUserCol := custommongo.NewMongoCollection(usercol)
	indexModels := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "email", Value: 1}, {Key: "username", Value: 1}}, // index in ascending order
			Options: options.Index().SetUnique(true),                               // make index unique
		},
	}

	_, err := customUserCol.CreateIndexes(context.Background(), indexModels)
	if err != nil {
		panic(err)
	}
	customRefreshTokenCol := custommongo.NewMongoCollection(refreshtokencol)

	BlogCollections = Domain.BlogCollections{
		Users:         customUserCol,
		RefreshTokens: customRefreshTokenCol,
	}
	// Initialize the router

	Router = gin.Default()

	// user router
	UserRouter()

	return Router

}
