package internal

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"
)

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
	err = json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		log.Printf("JSON decode error: %v", err)
		http.Error(w, "Invalid request body format", http.StatusBadRequest)
		return
	}

	username := loginRequest.Username
	password := loginRequest.Password

	user, ok := users[username]
	if !ok || !checkPasswordHash(password, user.HashedPassword) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	sessionToken := generateToken(32)
	csrfToken := generateToken(32)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    sessionToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    csrfToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HttpOnly: false,
		Secure:   true,
	})

	user.SessionToken = sessionToken
	user.CSRFToken = csrfToken
	users[username] = user

	response := map[string]string{"message": "Login successful!"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
