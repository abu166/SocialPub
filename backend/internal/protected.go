package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Protected handles requests to a protected resource, validating CSRF and session tokens.
func Protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Extract the session token from cookies
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Session token missing", http.StatusUnauthorized)
		return
	}

	sessionToken := sessionCookie.Value

	// Validate session token and retrieve the user
	var currentUser *NewLogin
	for _, user := range users {
		if user.SessionToken == sessionToken {
			currentUser = user // `user` is now already a pointer, no need for &user
			break
		}
	}

	if currentUser == nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	// Extract and validate the CSRF token from headers
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" || csrfToken != currentUser.CSRFToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	// Send user information as a response
	response := map[string]interface{}{
		"message":  fmt.Sprintf("CSRF validation successful! Welcome, %s", currentUser.Username),
		"is_admin": currentUser.IsAdmin,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
