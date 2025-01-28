package internal

import (
	"encoding/json"
	"net/http"
)

type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RegisterResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

func Register(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	username := req.Username
	email := req.Email
	password := req.Password

	// Validate input
	if username == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Message: "Username, email, and password are required.",
		})
		return
	}

	// Validate email format
	if !isValidEmail(email) {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Message: "Invalid email format.",
		})
		return
	}

	// Check if user already exists
	if _, exists := users[username]; exists {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(RegisterResponse{
			Success: false,
			Message: "Username already exists.",
		})
		return
	}

	// Hash password and save the user
	hashedPassword, err := hashPassword(password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Save user details
	users[username] = NewLogin{
		HashedPassword: hashedPassword,
		SessionToken:   "",
		CSRFToken:      "",
	}

	// Return success response
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: true,
		Message: "Registration successful.",
	})
}
