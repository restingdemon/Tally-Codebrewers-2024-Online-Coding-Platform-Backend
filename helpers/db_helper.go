package helpers

import (
	"context"
	"fmt"
	"worldwide-coders/models"
	"worldwide-coders/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

// Helper_InsertProblem inserts a new problem into the database and automatically increments the pid.
func Helper_InsertProblem(problem *models.Problem) (*mongo.InsertOneResult, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("problems")

	// Set options to sort by pid in descending order
	findOptions := options.FindOne().SetSort(bson.D{{"pid", -1}})

	// Find the problem with the highest pid
	var lastProblem models.Problem
	err := collection.FindOne(context.Background(), bson.M{}, findOptions).Decode(&lastProblem)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, fmt.Errorf("failed to find last problem: %s", err)
	}

	// Increment pid
	if lastProblem.Pid != 0 {
		problem.Pid = lastProblem.Pid + 1
	} else {
		problem.Pid = 1 // First problem
	}

	// Insert the new problem
	result, err := collection.InsertOne(context.Background(), problem)
	if err != nil {
		return nil, fmt.Errorf("failed to insert problem: %s", err)
	}

	return result, nil
}

func Helper_GetProblemByID(id int32) (*models.Problem, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("problems")
	problem := &models.Problem{}
	err := collection.FindOne(context.Background(), bson.M{"pid": id}).Decode(&problem)
	return problem, err
}

func Helper_GetAllProblems() ([]models.Problem, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("problems")
	cursor, err := collection.Find(context.Background(), bson.M{"visibility": true})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.Background())

	var problems []models.Problem
	if err := cursor.All(context.Background(), &problems); err != nil {
		return nil, err
	}

	return problems, nil
}

func Helper_GetNotVisibleProblems(role string, email string) ([]models.Problem, error) {
	collection := models.DB.Database("WorldwideCodersDb").Collection("problems")
	if role == utils.SuperAdminRole {
		cursor, err := collection.Find(context.Background(), bson.M{"visibility": false})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		var problems []models.Problem
		if err := cursor.All(context.Background(), &problems); err != nil {
			return nil, err
		}

		return problems, nil
	} else {
		cursor, err := collection.Find(context.Background(), bson.M{"author_id": email, "visibility": false})
		if err != nil {
			return nil, err
		}
		defer cursor.Close(context.Background())

		var problems []models.Problem
		if err := cursor.All(context.Background(), &problems); err != nil {
			return nil, err
		}

		return problems, nil
	}
}

func Helper_UpdateProblem(id int32, problem *models.Problem) error {
	collection := models.DB.Database("WorldwideCodersDb").Collection("problems")

	// Update the problem in the database
	_, err := collection.UpdateOne(
		context.Background(),
		bson.M{"pid": id},
		bson.M{
			"$set": bson.M{
				"title":       problem.Title,
				"description": problem.Description,
				"constraints": problem.Constraints,
				"test_cases":  problem.TestCases,
				"author_id":   problem.AuthorID,
				"visibility":  problem.Visibility,
			},
		},
	)
	return err
}
