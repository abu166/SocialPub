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

func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
	// Set response header for JSON content
	w.Header().Set("Content-Type", "application/json")

	// Hardcoded CSRF token (for now)
	csrfToken := "zcj7gYl_VoMa8pxJM78tLGRzpRJQCSpksh9F41hf-Fc"

	// Create JSON response
	response := map[string]string{"csrf_token": csrfToken}

	// Write JSON response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

//func GetCSRFToken(w http.ResponseWriter, r *http.Request) {
//	// Ensure the request is a GET
//	if r.Method != http.MethodGet {
//		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
//		return
//	}
//
//	// Retrieve the CSRF token from the cookie
//	csrfCookie, err := r.Cookie("X-CSRF-Token")
//	if err != nil || csrfCookie.Value == "" {
//		log.Printf("CSRF token retrieval failed: %v", err)
//		http.Error(w, "CSRF token not found", http.StatusUnauthorized)
//		return
//	}
//
//	// Send the CSRF token as a response
//	response := map[string]string{
//		"csrf_token": csrfCookie.Value,
//	}
//	w.Header().Set("Content-Type", "application/json")
//	w.WriteHeader(http.StatusOK)
//	if err := json.NewEncoder(w).Encode(response); err != nil {
//		log.Printf("Failed to encode CSRF token response: %v", err)
//	}
//}
