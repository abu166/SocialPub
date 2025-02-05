package internal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"main/internal/email"
	"net/http"
	"net/http/httptest"
	"testing"
)

// MockEmailSender simulates email sending during tests
func MockEmailSender(_, _ string) error {
	return nil
}

func TestRegisterAndLogin(t *testing.T) {
	fmt.Println("Starting TestRegisterAndLogin...")

	// Mock email sender
	EmailSender = MockEmailSender
	defer func() { EmailSender = email.SendEmail }() // Reset after test

	// Reset users
	users = make(map[string]*NewLogin)

	// Test Registration
	fmt.Println("Testing user registration...")
	registerReq := RegisterRequest{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "testpass123",
	}
	body, _ := json.Marshal(registerReq)

	rr := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/register", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	Register(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Registration failed: expected status %d, got %d", http.StatusOK, rr.Code)
	}
	fmt.Println("Registration successful!")

	// Test Login
	fmt.Println("Testing user login...")
	loginReq := LoginRequest{
		Username: "testuser",
		Password: "testpass123",
	}
	body, _ = json.Marshal(loginReq)

	rr = httptest.NewRecorder()
	req = httptest.NewRequest("POST", "/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	Login(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("Login failed: expected status %d, got %d", http.StatusOK, rr.Code)
	}
	fmt.Println("Login successful!")

	// Verify response
	var response map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&response); err != nil {
		t.Fatal("Failed to decode response:", err)
	}

	if response["csrf_token"] == "" || response["is_logged_in"] != true {
		t.Errorf("Invalid login response: %+v", response)
	}

	// Verify cookies
	if findCookie(rr.Result().Cookies(), "session_token") == nil {
		t.Error("Session cookie not set")
	}

	if findCookie(rr.Result().Cookies(), "csrf_token") == nil {
		t.Error("CSRF cookie not set")
	}

	fmt.Println("TestRegisterAndLogin completed successfully!")
}

func findCookie(cookies []*http.Cookie, name string) *http.Cookie {
	for _, c := range cookies {
		if c.Name == name {
			return c
		}
	}
	return nil
}
