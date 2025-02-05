package internal

import (
	"encoding/base64"
	"fmt"
	"testing"
)

// TestGenerateToken tests the generateToken function
func TestGenerateToken(t *testing.T) {
	// Test case 1: Check if the length of the generated token matches the expected length
	tokenLength := 32
	fmt.Println("Step 1: Generating a token of length", tokenLength)
	token := generateToken(tokenLength)
	fmt.Println("Generated token:", token)

	// Check if the generated token has the correct length
	expectedLength := 44 // Base64 encoded string will have a length of 4 * ceil(length/3)
	actualLength := len(token)
	fmt.Printf("Step 2: Checking token length. Expected: %d, Got: %d\n", expectedLength, actualLength)

	if actualLength != expectedLength {
		t.Errorf("Expected token length of %d, but got %d", expectedLength, actualLength)
	} else {
		fmt.Println("Step 2: Token length is correct.")
	}

	// Test case 2: Check if the token is a valid base64 URL-encoded string
	fmt.Println("Step 3: Checking if the token is a valid base64 URL-encoded string.")
	_, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		t.Errorf("Step 3: Token is not a valid base64 URL-encoded string: %v", err)
	} else {
		fmt.Println("Step 3: Token is a valid base64 URL-encoded string.")
	}

	// Test case 3: Ensure the function does not log an error for a valid length
	// This is indirectly tested since the function will panic if there's an error
	// If the test passes, it means the function worked without issues
	fmt.Println("Step 4: Function executed without errors (if you see no panic).")
	t.Log("Step 4: Token generated successfully")
}
