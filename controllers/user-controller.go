// auth-controller.go
package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"worldwide-coders/helpers"
	"worldwide-coders/models"
	"worldwide-coders/utils"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type GoogleUser struct {
	ID          primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Email       string             `json:"email" bson:"email"`
	Name        string             `json:"name" bson:"name"`
	Phone       string             `json:"phone" bson:"phone"`
	Description string             `json:"description" bson:"description"`
	Token       string             `json:"token"`
	Image       string             `json:"image" bson:"image"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the token
	var user = &GoogleUser{}
	utils.ParseBody(r, user)

	if user.Email == "" {
		http.Error(w, fmt.Sprintf("No email provided"), http.StatusBadRequest)
		return
	}

	if !utils.IsloginValid(user.Email, user.Token) {
		http.Error(w, fmt.Sprintf("User token not valid"), http.StatusBadRequest)
		return
	}
	// Check if the user already exists in the database based on their email
	existingUser, err := helpers.Helper_GetUserByEmail(user.Email)
	if err != nil && !errors.Is(err, mongo.ErrNoDocuments) {
		http.Error(w, fmt.Sprintf("Failed to check user existence: %s", err), http.StatusInternalServerError)
		return
	}

	// If the user doesn't exist, create a new user in the database
	if existingUser == nil {
		err := createUser(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to create user: %s", err), http.StatusInternalServerError)
			return
		}
		existingUser, err = helpers.Helper_GetUserByEmail(user.Email)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to find user after create: %s", err), http.StatusInternalServerError)
		}
	} else {
		existingUser.Image = user.Image
		err = helpers.Helper_UpdateUser(existingUser)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to update user: %s", err), http.StatusInternalServerError)
			return
		}
	}

	// Generate JWT tokens for the user
	token, refreshToken, err := helpers.GenerateAllTokens(existingUser.Email, existingUser.Name, existingUser.Role, existingUser.ID.Hex())
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate tokens: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the user data and tokens as a JSON response
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"_id":         existingUser.ID.Hex(),
			"email":       existingUser.Email,
			"name":        existingUser.Name,
			"phone":       existingUser.Phone,
			"descriptiom": existingUser.Description,
			"role":        existingUser.Role,
			"image":       existingUser.Image,
		},
		"token":         token,
		"refresh_token": refreshToken,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

}

func createUser(googleUser *GoogleUser) error {
	collection := models.DB.Database("WorldwideCodersDb").Collection("users")
	user := models.User{
		Email: googleUser.Email,
		Name:  googleUser.Name,
		Role:  utils.UserRole,
		Image: googleUser.Image,
	}
	if googleUser.Email == "akshay.garg130803@gmail.com" {
		user.Role = utils.SuperAdminRole
	}
	_, err := collection.InsertOne(context.Background(), user)
	if err != nil {
		return fmt.Errorf("failed to insert user: %s", err)
	}

	return nil
}

func GetUserByEmail(w http.ResponseWriter, r *http.Request) {
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	if email == "" {
		users, err := helpers.Helper_ListAllUsers()
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to list all users: %s", err), http.StatusInternalServerError)
			return
		}

		response, err := json.Marshal(users)
		if err != nil {
			http.Error(w, "Failed to marshal users", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	} else {
		user, err := helpers.Helper_GetUserByEmail(email)
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				http.Error(w, "User not found", http.StatusNotFound)
			} else {
				http.Error(w, fmt.Sprintf("Failed to get User: %s", err), http.StatusInternalServerError)
			}
			return
		}

		response, err := json.Marshal(user)
		if err != nil {
			http.Error(w, "Failed to marshal user details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Parse request body to get the updated user details
	var updatedUser = &models.User{}
	utils.ParseBody(r, updatedUser)

	// Extract email from the context
	emailValue := r.Context().Value("email")
	if emailValue == nil {
		http.Error(w, "Email not found in context", http.StatusInternalServerError)
		return
	}

	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}

	// Retrieve the user from the database based on the email
	existingUser, err := helpers.Helper_GetUserByEmail(email)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, fmt.Sprintf("User not found"), http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get user: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Convert existingUser to GoogleUser for the update process
	updatedUser = &models.User{
		ID:          existingUser.ID,
		Email:       existingUser.Email,
		Name:        existingUser.Name,
		Phone:       updatedUser.Phone,
		Description: updatedUser.Description,
		Role:        existingUser.Role,
		Image:       existingUser.Image,
	}

	// Update the user in the database
	err = helpers.Helper_UpdateUser(updatedUser)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to update user: %s", err), http.StatusInternalServerError)
		return
	}

	// Return the updated user data as a JSON response
	response := map[string]interface{}{
		"user": map[string]interface{}{
			"_id":         updatedUser.ID.Hex(),
			"email":       updatedUser.Email,
			"name":        updatedUser.Name,
			"phone":       updatedUser.Phone,
			"description": updatedUser.Description,
			"role":        updatedUser.Role,
			"image":       updatedUser.Image,
		},
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
