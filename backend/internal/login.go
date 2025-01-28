package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

// Store user tokens in memory for demonstration purposes (not recommended for production)
var userTokens = map[string]string{}

// Login handles user authentication and issues session & CSRF tokens
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	var loginRequest LoginRequest

	// Read and log the raw body for debugging
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}
	log.Printf("Raw request body: %s", string(body))

	// Reset the body and decode JSON
	r.Body = io.NopCloser(bytes.NewBuffer(body))
	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	username := loginRequest.Username
	password := loginRequest.Password

	// Validate user credentials
	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate session and CSRF tokens
	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	// Store CSRF token for the user
	userTokens[username] = csrfToken

	// Update user record with tokens
	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	// Set cookies
	setCookie(w, "session_token", sessionToken, true, true)
	setCookie(w, "csrf_token", csrfToken, true, false)

	// Respond with CSRF token
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"csrf_token": csrfToken}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

// Helper function to set cookies
func setCookie(w http.ResponseWriter, name, value string, httpOnly, secure bool) {
	http.SetCookie(w, &http.Cookie{
		Name:     name,
		Value:    value,
		Path:     "/",
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: httpOnly,
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
	})
}
