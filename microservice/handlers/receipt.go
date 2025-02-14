package handlers

import (
	"encoding/json"
	"log"
	"microservice/database"
	"microservice/services"
	"microservice/utils"
	"net/http"
)

func GenerateReceipt(w http.ResponseWriter, r *http.Request) {
	// Handle preflight requests
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Extract transaction ID from query parameters
	vars := r.URL.Query()
	transactionID := vars.Get("transactionID") // transactionID is already a string
	if transactionID == "" {
		log.Println("Error: Missing transaction ID")
		http.Error(w, "Missing transaction ID", http.StatusBadRequest)
		return
	}
	log.Printf("Processing receipt for transaction ID: %s", transactionID)

	// Mock transaction details
	details := map[string]string{
		"name": "John Doe",
	}

	// Generate the receipt PDF
	filePath, err := services.GenerateReceipt(transactionID, details)
	if err != nil {
		log.Printf("Error generating receipt for ID %s: %v", transactionID, err)
		http.Error(w, "Failed to generate receipt", http.StatusInternalServerError)
		return
	}
	log.Printf("PDF generated successfully at: %s", filePath)

	// Send the receipt via email
	email := "abulljk@gmail.com"
	subject := "Your Receipt"
	body := "Thank you for your purchase! Please find your receipt attached."
	err = utils.SendEmail(email, subject, body, filePath)
	if err != nil {
		log.Printf("Error sending email for ID %s: %v", transactionID, err)
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}
	log.Println("Email sent successfully")

	// Update the transaction status to "Completed"
	newStatus := "Completed"
	err = database.UpdateTransactionStatus(transactionID, newStatus)
	if err != nil {
		log.Printf("Error updating transaction status for ID %s: %v", transactionID, err)
		http.Error(w, "Failed to update transaction status", http.StatusInternalServerError)
		return
	}
	log.Println("Transaction status updated to Completed")

	// Respond with success and include the transaction status
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message":       "Receipt generated, sent, and transaction finalized successfully",
		"transactionID": transactionID,
		"status":        newStatus,
	})
}
