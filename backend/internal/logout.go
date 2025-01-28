package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

// Logout handles user logout by invalidating session and CSRF tokens
func Logout(w http.ResponseWriter, r *http.Request) {
	// Ensure the request is a POST
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the session token from cookies
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Unauthorized: No session token provided", http.StatusUnauthorized)
		return
	}

	sessionToken := sessionCookie.Value

	// Find the user associated with the session token
	var username string
	var currentUser *NewLogin
	for u, user := range users {
		if user.SessionToken == sessionToken {
			username = u
			currentUser = &user
			break
		}
	}

	if currentUser == nil {
		http.Error(w, "Unauthorized: Invalid session token", http.StatusUnauthorized)
		return
	}

	// Clear session and CSRF tokens from the user record
	currentUser.SessionToken = ""
	currentUser.CSRFToken = ""
	users[username] = *currentUser // Update the user in the map

	// Expire the session and CSRF cookies
	http.SetCookie(w, &http.Cookie{
		Name:     "session_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: true,
		Secure:   true,
	})

	http.SetCookie(w, &http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HttpOnly: false,
		Secure:   true,
	})

	// Respond to the client
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"message": "Logout successful"})
}
