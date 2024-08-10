package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"worldwide-coders/helpers"
	"worldwide-coders/models"
	"worldwide-coders/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

func CreateProblem(w http.ResponseWriter, r *http.Request) {
	var problem models.Problem
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	roleVal := r.Context().Value("role")
	role, ok1 := roleVal.(string)
	if !ok1 {
		http.Error(w, "Failed to retrieve role from context", http.StatusInternalServerError)
		return
	}
	problem.AuthorID = email
	if role == utils.UserRole {
		problem.Visibility = false
	} else {
		problem.Visibility = true
	}

	result, err := helpers.Helper_InsertProblem(&problem)
	if err != nil {
		http.Error(w, "Failed to create problem", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		http.Error(w, "Failed to marshal problem details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetProblems(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	id := queryParams.Get("id")
	// Fetch a specific problem by ID
	if id != "" {
		pid, err := strconv.Atoi(id)
		if err != nil {
			http.Error(w, "Invalid problem ID", http.StatusBadRequest)
			return
		}
		problem, err := helpers.Helper_GetProblemByID(int32(pid))
		if err != nil {
			if err == mongo.ErrNoDocuments {
				http.Error(w, "Problem not found", http.StatusNotFound)
			} else {
				http.Error(w, "Failed to fetch problem", http.StatusInternalServerError)
			}
			return
		}
		if len(problem.TestCases) > 2 {
			problem.TestCases = problem.TestCases[:2]
		}
		response, err := json.Marshal(problem)
		if err != nil {
			http.Error(w, "Failed to marshal problem details", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response)
		return
	}

	// Fetch all problems
	problems, err := helpers.Helper_GetAllProblems()
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Problem not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch problem", http.StatusInternalServerError)
		}
		return
	}
	for i := range problems {
		if len(problems[i].TestCases) > 2 {
			problems[i].TestCases = problems[i].TestCases[:2]
		}
	}
	response, err := json.Marshal(problems)
	if err != nil {
		http.Error(w, "Failed to marshal problem details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func GetNotVisibleProblems(w http.ResponseWriter, r *http.Request) {
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	roleVal := r.Context().Value("role")
	role, ok1 := roleVal.(string)
	if !ok1 {
		http.Error(w, "Failed to retrieve role from context", http.StatusInternalServerError)
		return
	}

	problems, err := helpers.Helper_GetNotVisibleProblems(role, email)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Problem not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch problem", http.StatusInternalServerError)
		}
		return
	}
	for i := range problems {
		if len(problems[i].TestCases) > 2 {
			problems[i].TestCases = problems[i].TestCases[:2]
		}
	}
	response, err := json.Marshal(problems)
	if err != nil {
		http.Error(w, "Failed to marshal problem details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func UpdateProblem(w http.ResponseWriter, r *http.Request) {
	emailValue := r.Context().Value("email")
	email, ok := emailValue.(string)
	if !ok {
		http.Error(w, "Failed to retrieve email from context", http.StatusInternalServerError)
		return
	}
	roleVal := r.Context().Value("role")
	role, ok1 := roleVal.(string)
	if !ok1 {
		http.Error(w, "Failed to retrieve role from context", http.StatusInternalServerError)
		return
	}

	pidStr := mux.Vars(r)["pid"]
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		http.Error(w, "Invalid problem ID", http.StatusBadRequest)
		return
	}

	existingproblem, err := helpers.Helper_GetProblemByID(int32(pid))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Problem not found", http.StatusNotFound)
		} else {
			http.Error(w, "Failed to fetch problem", http.StatusInternalServerError)
		}
		return
	}

	// Check if the user is the author or a superadmin
	if role == utils.UserRole && existingproblem.AuthorID != email {
		http.Error(w, "You can only update your own problems", http.StatusForbidden)
		return
	}
	problem := &models.Problem{}
	if err := json.NewDecoder(r.Body).Decode(&problem); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}
	if problem.Constraints != "" {
		existingproblem.Constraints = problem.Constraints
	}
	if problem.Description != "" {
		existingproblem.Description = problem.Description
	}
	if problem.TestCases != nil {
		existingproblem.TestCases = problem.TestCases
	}
	if problem.Title != "" {
		existingproblem.Title = problem.Title
	}
	if role == utils.UserRole {
		existingproblem.Visibility = false
	}
	if role == utils.SuperAdminRole{
		existingproblem.Visibility=problem.Visibility
	}
	// Update the problem in the database
	err = helpers.Helper_UpdateProblem(int32(pid), existingproblem)
	if err != nil {
		http.Error(w, "Failed to update problem", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(existingproblem)
}
