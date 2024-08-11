package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
	"worldwide-coders/helpers"
	"worldwide-coders/models"
	"worldwide-coders/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateContest(w http.ResponseWriter, r *http.Request) {
	var contest models.Contest
	if err := json.NewDecoder(r.Body).Decode(&contest); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	contest.HostID = email

	collection := models.DB.Database("WorldwideCodersDb").Collection("contests")
	if _, err := collection.InsertOne(context.Background(), contest); err != nil {
		http.Error(w, "Failed to create contest", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contest)
}

func ContestRegister(w http.ResponseWriter, r *http.Request) {
	contestId, err := primitive.ObjectIDFromHex(mux.Vars(r)["contestId"])
	if err != nil {
		http.Error(w, "Invalid contest ID", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value("email").(string)

	participant := models.Participant{
		ContestID: contestId,
		UserID:    userID,
		Score:     0,
	}

	collection := models.DB.Database("WorldwideCodersDb").Collection("participants")
	if _, err := collection.InsertOne(context.Background(), participant); err != nil {
		http.Error(w, "Failed to register for contest", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func GetContest(w http.ResponseWriter, r *http.Request) {
	contestId := r.URL.Query().Get("id")
	if contestId != "" {
		objectID, err := primitive.ObjectIDFromHex(contestId)
		if err != nil {
			http.Error(w, "Invalid contest ID", http.StatusBadRequest)
			return
		}

		contest, err := helpers.Helper_GetContestById(objectID)
		if err != nil {
			http.Error(w, "Contest not found", http.StatusNotFound)
			return
		}

		// Check if the current time is >= contest start time
		if time.Now().Unix() >= contest.StartTime {
			response, err := json.Marshal(contest)
			if err != nil {
				http.Error(w, "Failed to marshal contest details", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write(response)
		} else {
			http.Error(w, "Contest has not started yet", http.StatusForbidden)
		}
		return
	}

	// Fetch all contests without showing problem statements
	contests, err := helpers.Helper_GetAllContests()
	if err != nil {
		http.Error(w, "Failed to fetch contests", http.StatusInternalServerError)
		return
	}

	// Remove problem statements from each contest
	for i := range contests {
		contests[i].Problems = nil
	}

	response, err := json.Marshal(contests)
	if err != nil {
		http.Error(w, "Failed to marshal contests", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
func GetAllRegistrations(w http.ResponseWriter, r *http.Request) {
	contestId, err := primitive.ObjectIDFromHex(mux.Vars(r)["contestId"])
	if err != nil {
		http.Error(w, "Invalid contest ID", http.StatusBadRequest)
		return
	}

	role := r.Context().Value("role").(string)
	hostID := r.Context().Value("email").(string)

	var filter bson.M
	if role == utils.UserRole {
		contest, err := helpers.Helper_GetContestById(contestId)
		if err != nil {
			http.Error(w, "Contest not found", http.StatusNotFound)
			return
		}
		if contest.HostID != hostID {
			http.Error(w, "Not authorised to view contest participants", http.StatusNotFound)
			return
		}
	}
	filter = bson.M{"contest_id": contestId}

	collection := models.DB.Database("WorldwideCodersDb").Collection("participants")
	cursor, err := collection.Find(context.Background(), filter)
	if err != nil {
		http.Error(w, "Failed to get registrations", http.StatusInternalServerError)
		return
	}

	var participants []models.Participant
	if err := cursor.All(context.Background(), &participants); err != nil {
		http.Error(w, "Failed to decode participants", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(participants)
}
func CheckRegistration(w http.ResponseWriter, r *http.Request) {
	// Extract contestId from URL path
	contestId, err := primitive.ObjectIDFromHex(mux.Vars(r)["contestId"])
	if err != nil {
		http.Error(w, "Invalid contest ID", http.StatusBadRequest)
		return
	}

	// Retrieve email from context
	email, ok := r.Context().Value("email").(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}

	// Check if the contest exists
	existingContest, err := helpers.Helper_GetContestById(contestId)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			http.Error(w, "Contest not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("Failed to get contest: %s", err), http.StatusInternalServerError)
		}
		return
	}

	// Define response structure
	type Response struct {
		IsRegistered         bool                `json:"is_registered"`
		ExistingRegistration *models.Participant `json:"existing_registration,omitempty"`
	}

	response := &Response{}

	// Check if the user is already registered for the contest
	existingRegistration, err := helpers.Helper_GetRegistrationByEmailAndContest(email, existingContest.ContestID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to check registration: %v", err), http.StatusInternalServerError)
		return
	}

	// If registration exists, set response accordingly
	if existingRegistration != nil {
		response.IsRegistered = true
		response.ExistingRegistration = existingRegistration
	} else {
		response.IsRegistered = false
	}

	// Marshal response to JSON and send it back
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}
func GetLeaderboard(w http.ResponseWriter, r *http.Request) {
	contestId, err := primitive.ObjectIDFromHex(mux.Vars(r)["contestId"])
	if err != nil {
		http.Error(w, "Invalid contest ID", http.StatusBadRequest)
		return
	}

	var leaderboard models.Leaderboard
	collection := models.DB.Database("WorldwideCodersDb").Collection("leaderboards")
	err = collection.FindOne(context.Background(), bson.M{"contest_id": contestId}).Decode(&leaderboard)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Leaderboard not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch leaderboard", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(leaderboard)
}
