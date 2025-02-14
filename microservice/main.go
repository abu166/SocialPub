package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors" // Import the CORS middleware library
	"microservice/database"
	"microservice/handlers"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials from environment variables
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	dbSSLMode := os.Getenv("DB_SSLMODE")

	// Construct the data source name (DSN)
	dataSourceName := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		dbHost, dbUser, dbPassword, dbName, dbPort, dbSSLMode)

	// Initialize the database
	database.InitDB(dataSourceName)

	// Create a new router
	r := mux.NewRouter()

	// Define routes
	r.HandleFunc("/payment", handlers.HandlePayment).Methods("POST")
	r.HandleFunc("/receipt", handlers.GenerateReceipt).Methods("GET")

	// Add CORS middleware
	corsHandler := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Allow requests from frontend
		AllowedMethods: []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type", "X-CSRF-Token", "Authorization"},
	})

	// Wrap the router with the CORS middleware
	handler := corsHandler.Handler(r)

	// Start the server
	log.Println("Microservice running on port 8081...")
	log.Fatal(http.ListenAndServe(":8081", handler))
}
