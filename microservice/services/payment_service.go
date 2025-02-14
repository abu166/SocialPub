package services

import (
	"microservice/models"
	"time"
)

func ValidatePayment(paymentForm models.PaymentForm) bool { // Use the shared model
	// Mock validation logic
	expirationDate := parseDate(paymentForm.ExpirationDate)
	if expirationDate.Before(time.Now()) {
		return false // Card expired
	}

	// Simulate successful payment
	return true
}

func parseDate(dateStr string) time.Time {
	layout := "01/06"
	date, _ := time.Parse(layout, dateStr)
	return date
}
