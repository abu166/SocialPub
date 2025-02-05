package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"testing"
)

func TestEnvFileLoading(t *testing.T) {
	fmt.Println("Starting TestEnvFileLoading...")

	envPath := "../.env"
	if err := godotenv.Load(envPath); err != nil {
		t.Fatalf("Failed to load .env file from %s: %v", envPath, err)
	}

	expectedKey := "SERVER_PORT"
	value, exists := os.LookupEnv(expectedKey)
	if !exists || value == "" {
		t.Errorf("Expected environment variable %s to be set, but got empty or missing value", expectedKey)
	}

	fmt.Println("TestEnvFileLoading completed successfully!")
}
