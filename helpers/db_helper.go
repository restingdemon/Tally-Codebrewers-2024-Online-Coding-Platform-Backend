package helpers

import (
	"context"
	"fmt"
	"worldwide-coders/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// **********USER************************

func Helper_GetUserByID(userID primitive.ObjectID) (*models.User, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("users")

	filter := bson.M{"_id": userID}
	user := &models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Helper_GetUserByEmail(email string) (*models.User, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("users")

	filter := bson.M{"email": email}
	user := &models.User{}
	err := collection.FindOne(context.TODO(), filter).Decode(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func Helper_ListAllUsers() ([]models.User, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("users")

	cursor, err := collection.Find(context.TODO(), bson.D{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var users []models.User
	if err := cursor.All(context.TODO(), &users); err != nil {
		return nil, fmt.Errorf("failed to decode users: %s", err)
	}

	return users, nil
}

func Helper_UpdateUser(user *models.User) error {
	collection := models.DB.Database("WorldwideCodersDb").Collection("users")

	update := bson.M{
		"$set": models.User{
			Email:       user.Email,
			Name:        user.Name,
			Phone:       user.Phone,
			Description: user.Description,
			Role:        user.Role,
			Image:       user.Image,
		},
	}

	// Update user in the database based on the email
	_, err := collection.UpdateOne(context.Background(), bson.M{"email": user.Email}, update)
	if err != nil {
		return fmt.Errorf("failed to update user: %s", err)
	}

	return nil
}
