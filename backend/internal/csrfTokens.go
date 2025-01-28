package internal

import (
	"encoding/base64"
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
)

func generateToken(length int) string {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("Failed to generate token: %v", err)
	}

	return base64.URLEncoding.EncodeToString(bytes) // Base64 converts bytes into a string

}

//func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
//	// Set response header for JSON content
//	w.Header().Set("Content-Type", "application/json")
//
//	// Hardcoded CSRF token (for now)
//	csrfToken := "zcj7gYl_VoMa8pxJM78tLGRzpRJQCSpksh9F41hf-Fc"
//
//	// Create JSON response
//	response := map[string]string{"csrf_token": csrfToken}
//
//	// Write JSON response
//	if err := json.NewEncoder(w).Encode(response); err != nil {
//		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
//		return
//	}
//}

// GetCSRFToken returns the CSRF token stored in the user's cookie
func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil || sessionCookie.Value == "" {
		http.Error(w, "Unauthorized: missing session token", http.StatusUnauthorized)
		return
	}

	// Validate session token and get associated user
	var csrfToken string
	for _, u := range users {
		if u.SessionToken == sessionCookie.Value {
			csrfToken = u.CSRFToken
			break
		}
	}

	if csrfToken == "" {
		http.Error(w, "Unauthorized: invalid session token", http.StatusUnauthorized)
		return
	}

	// Respond with CSRF token
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"csrf_token": csrfToken}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func HandleCsrfToken(w http.ResponseWriter, r *http.Request) {
	sessionCookie, err := r.Cookie("session_token")
	if err != nil {
		http.Error(w, "Unauthorized: missing session token", http.StatusUnauthorized)
		return
	}

	// Validate session token
	var user *NewLogin
	for _, u := range users {
		if u.SessionToken == sessionCookie.Value {
			user = &u
			break
		}
	}

	if user == nil {
		http.Error(w, "Unauthorized: invalid session token", http.StatusUnauthorized)
		return
	}

	// Return CSRF token
	w.Header().Set("Content-Type", "application/json")
	response := map[string]string{"csrf_token": user.CSRFToken}
	json.NewEncoder(w).Encode(response)
}
