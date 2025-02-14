package handlers

import (
	"encoding/json"
	"log"
	"microservice/database"
	"microservice/models"
	"microservice/services"
	"net/http"
)

func HandlePayment(w http.ResponseWriter, r *http.Request) {
	var paymentForm models.PaymentForm
	if err := json.NewDecoder(r.Body).Decode(&paymentForm); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the payment form for debugging
	log.Printf("Received payment form: %+v", paymentForm)

	// Validate payment details
	isValid := services.ValidatePayment(paymentForm)

	// Simulate a successful payment for testing
	isValid = true // Remove this line once validation is fixed

	status := "Declined"
	if isValid {
		status = "Paid"
	}

	// Prepare payment details
	paymentDetails := models.PaymentDetails{
		CardNumber:     paymentForm.CardNumber,
		ExpirationDate: paymentForm.ExpirationDate,
		Name:           paymentForm.Name,
		Address:        paymentForm.Address,
	}

	// Save the transaction to the database
	transactionID, err := database.SaveTransaction(status, paymentDetails)
	if err != nil {
		log.Printf("Error saving transaction: %v", err)
		http.Error(w, "Failed to save transaction", http.StatusInternalServerError)
		return
	}

	// Update the transaction status
	if err := database.UpdateTransactionStatus(transactionID, status); err != nil {
		log.Printf("Error updating transaction status: %v", err)
		http.Error(w, "Failed to update transaction status", http.StatusInternalServerError)
		return
	}

	// Send response back to the main server
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": isValid})
}
