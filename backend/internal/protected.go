package internal

import (
	"fmt"
	"net/http"
)

// Protected handles requests to a protected resource, validating CSRF and session tokens.
func Protected(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
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
			currentUser = &user
			break
		}
	}

	if currentUser == nil {
		http.Error(w, "Invalid session token", http.StatusUnauthorized)
		return
	}

	// Extract and validate the CSRF token from headers
	csrfToken := r.Header.Get("X-CSRF-Token")
	if csrfToken == "" {
		http.Error(w, "CSRF token missing", http.StatusForbidden)
		return
	}

	if csrfToken != currentUser.CSRFToken {
		http.Error(w, "Invalid CSRF token", http.StatusForbidden)
		return
	}

	// Retrieve username and send a success response
	username := r.FormValue("username")
	if username == "" {
		username = "Anonymous"
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "CSRF validation successful! Welcome, %s", username)
}
