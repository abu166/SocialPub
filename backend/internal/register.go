package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"main/internal/email"
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON input", http.StatusBadRequest)
		return
	}

	// Validate input
	if err := validateInput(req.Username, req.Email, req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if user already exists (assume `users` is a global map)
	if _, exists := users[req.Username]; exists {
		http.Error(w, "Username already exists.", http.StatusConflict)
		return
	}

	// Hash password
	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	// Save user details
	users[req.Username] = NewLogin{
		HashedPassword: hashedPassword,
		SessionToken:   "",
		CSRFToken:      "",
	}

	// Generate confirmation code
	rand.Seed(time.Now().UnixNano())
	confirmationCode := rand.Intn(900000) + 100000 // 6-digit code

	// Send email
	if err := email.SendEmail(req.Email, fmt.Sprintf("%d", confirmationCode)); err != nil {
		log.Println("Error sending email:", err)
		http.Error(w, "Failed to send confirmation email", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	json.NewEncoder(w).Encode(RegisterResponse{
		Success: true,
		Message: "Confirmation email sent! Check your inbox.",
	})
}

func validateInput(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("Username, email, and password are required.")
	}
	if len(password) < 6 {
		return errors.New("Password must be at least 6 characters long.")
	}
	return nil
}
