package internal

import (
	"encoding/json"
	"errors"
	"log"
	"main/internal/email"
	"net/http"
)

// Session store (temporary; use database for production)
var sessionStore = make(map[string]*NewLogin)

var EmailSender = email.SendEmail

// ValidateSessionAndCSRF checks session and CSRF tokens.
func ValidateSessionAndCSRF(r *http.Request) (*NewLogin, error) {
	// Extract session token from cookies
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		return nil, errors.New("session token missing")
	}

	sessionToken := sessionCookie.Value
	user, exists := sessionStore[sessionToken]
	if !exists {
		return nil, errors.New("invalid session token")
	}

	// Extract CSRF token from headers
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" || csrfToken != user.CSRFToken {
		return nil, errors.New("invalid CSRF token")
	}

	return user, nil
}

func CheckAuth(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		log.Println("No session cookie found")
		http.Error(w, "Unauthorized: missing session token", http.StatusUnauthorized)
		return
	}

	log.Println("Received session token:", sessionCookie.Value)

	var user *NewLogin
	for _, u := range users {
		if u.SessionToken == sessionCookie.Value {
			user = u
			break
		}
	}

	if user == nil {
		log.Println("Invalid session token:", sessionCookie.Value)
		http.Error(w, "Unauthorized: invalid session token", http.StatusUnauthorized)
		return
	}

	// Return authentication status
	response := map[string]bool{
		"is_logged_in": true,
		"is_admin":     user.IsAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
