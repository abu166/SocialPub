package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Print environment variables
	fmt.Println("SMTP_HOST:", os.Getenv("SMTP_HOST"))
	fmt.Println("SMTP_PORT:", os.Getenv("SMTP_PORT"))
	fmt.Println("SMTP_USER:", os.Getenv("SMTP_USER"))
	fmt.Println("SMTP_PASSWORD:", os.Getenv("SMTP_PASS"))
}
