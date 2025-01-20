package main

import (
	"log"
	"net/http"

	"main/internal/handlers"
	"main/internal/middleware"
	"main/pkg/database"
)

func main() {
	// Initialize the database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// Initialize routes
	router := http.NewServeMux()

	// Apply middleware and handlers
	router.HandleFunc("/post", middleware.RateLimiter(handlers.PostHandler))
	router.HandleFunc("/get", middleware.RateLimiter(handlers.GetHandler))
	router.HandleFunc("/users", middleware.RateLimiter(middleware.CORS(handlers.UserHandler)))
	router.HandleFunc("/user/create", middleware.RateLimiter(middleware.CORS(handlers.CreateUser)))
	router.HandleFunc("/user/update", middleware.RateLimiter(middleware.CORS(handlers.UpdateUserHandler)))
	router.HandleFunc("/user/delete", middleware.RateLimiter(middleware.CORS(handlers.DeleteUserHandler)))
	router.HandleFunc("/user/get", middleware.RateLimiter(middleware.CORS(handlers.GetUserHandler)))
	router.HandleFunc("/send-email", middleware.CORS(handlers.SendEmailHandler))

	// Serve static files
	fs := http.FileServer(http.Dir("./static"))
	router.Handle("/", fs)

	// Start server
	log.Println("Server is running on port 8080...")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
