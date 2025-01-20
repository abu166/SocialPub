package main

import (
	"log"
	"main/internal"
	"net/http"
)

var allowedOrigins = []string{"http://localhost:3000", "http://172.20.10.2:3000"}

func main() {
	mux := http.NewServeMux()

	// Register routes with middleware applied
	mux.Handle("/csrf-token", corsMiddleware(http.HandlerFunc(internal.GetCSRFToken)))
	mux.Handle("/register", corsMiddleware(http.HandlerFunc(internal.Register)))
	mux.Handle("/login", corsMiddleware(http.HandlerFunc(internal.Login)))
	mux.Handle("/logout", corsMiddleware(http.HandlerFunc(internal.Logout)))
	mux.Handle("/protected", corsMiddleware(http.HandlerFunc(internal.Protected)))

	// Serve static files under /static/
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Start server
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("Request Origin: %s", origin) // Debugging origin

		// Allow CORS for matching origins
		for _, o := range allowedOrigins {
			if origin == o {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				break
			}
		}

		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		// Include X-CSRF-Token in the allowed headers
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-CSRF-Token")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
