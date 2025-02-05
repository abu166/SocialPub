package internal

import (
	"fmt"
	"net/http"
)

// AuthorsRoute restricts access to authors only.
func AuthorsRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	// Validate session
	currentUser, err := ValidateSessionAndCSRF(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Ensure the user is an author
	if !currentUser.IsAdmin {
		http.Error(w, "Forbidden: Access restricted to authors", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Welcome to the authors' section, %s!", currentUser.Username)
}
