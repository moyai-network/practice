package data

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ctx returns a context.Context.
func ctx() context.Context {
	return context.Background()
}

// db is the Upper database session.
var db *mongo.Database

// init creates the Upper database connection.
func init() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://localhost"))
	if err != nil {
		panic(err)
	}
	db = client.Database("practice")

	userCollection = db.Collection("users")
}

const Salt = "McBurgerLOLSuperCoolPakistaniBeefMan"
