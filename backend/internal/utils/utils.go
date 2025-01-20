package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"main/internal/models"
	"net/http"
	"strings"

	"github.com/sirupsen/logrus"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidInput      = errors.New("invalid input")
	ErrDuplicateEmail    = errors.New("email already exists")
	ErrDatabaseOperation = errors.New("database operation failed")
)

func SendJSONResponse(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		Log.WithError(err).Error("Failed to encode JSON response")
	}
}

func SendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	SendJSONResponse(w, statusCode, models.ResponseData{
		Status:  "error",
		Message: message,
	})
}

func HandleError(w http.ResponseWriter, err error, statusCode int, logger *logrus.Entry) {
	logger.WithError(err).Error("Operation failed")

	var response models.ResponseData
	switch {
	case errors.Is(err, ErrUserNotFound):
		response = models.ResponseData{Status: "error", Message: "User not found"}
	case errors.Is(err, ErrInvalidInput):
		response = models.ResponseData{Status: "error", Message: "Invalid input provided"}
	case errors.Is(err, ErrDuplicateEmail):
		response = models.ResponseData{Status: "error", Message: "Email already exists"}
	default:
		response = models.ResponseData{Status: "error", Message: "Internal server error"}
	}

	SendJSONResponse(w, statusCode, response)
}

func ValidateUser(user models.User) error {
	if user.UserName == "" {
		return fmt.Errorf("%w: username is required", ErrInvalidInput)
	}
	if user.UserEmail == "" {
		return fmt.Errorf("%w: email is required", ErrInvalidInput)
	}
	if !IsValidEmail(user.UserEmail) {
		return fmt.Errorf("%w: invalid email format", ErrInvalidInput)
	}
	return nil
}

func IsValidEmail(email string) bool {
	if len(email) < 3 || len(email) > 254 {
		return false
	}

	atIndex := strings.Index(email, "@")
	dotIndex := strings.LastIndex(email, ".")

	return atIndex > 0 &&
		dotIndex > atIndex &&
		dotIndex < len(email)-1
}

func IsDuplicateEmailError(err error) bool {
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
