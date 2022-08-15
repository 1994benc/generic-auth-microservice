package users

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateUser(db *mongo.Client, user *User) (string, error) {
	result, err := db.Database(os.Getenv("DATABASE_NAME")).Collection("users").InsertOne(context.TODO(), user)
	return result.InsertedID.(primitive.ObjectID).Hex(), err
}
