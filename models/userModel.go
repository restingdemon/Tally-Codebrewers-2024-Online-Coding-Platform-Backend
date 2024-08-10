package models

import (
	"worldwide-coders/database"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var DB *mongo.Client

type User struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Name        string             `json:"name" bson:"name"`
	Phone       string             `json:"phone" bson:"phone"`
	Description string             `json:"description" bson:"description"`
	Image       string             `json:"image" bson:"image"`
	Role        string             `json:"role" bson:"role"`
}

func init() {
	database.Connect()
	DB = database.GetDB()
}
