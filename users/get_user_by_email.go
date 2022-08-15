package users

import (
	"context"
	"os"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

func GetUserByEmail(db *mongo.Client, email string) (*User, error) {
	filter := bson.D{{Key: "email", Value: email}}
	var user User
	if err := db.Database(os.Getenv("DATABASE_NAME")).Collection("users").FindOne(context.TODO(), filter).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
