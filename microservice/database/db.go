package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	_ "github.com/lib/pq" // PostgreSQL driver
	"log"
	"microservice/models"
)

var DB *sql.DB

func InitDB(dataSourceName string) {
	var err error
	DB, err = sql.Open("postgres", dataSourceName)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Ping to verify connection
	if err := DB.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Connected to the database!")
}

// SaveTransaction inserts a new transaction into the database.
func SaveTransaction(status string, paymentDetails models.PaymentDetails) (string, error) {
	// Generate a unique transaction ID as a string
	transactionID := uuid.New().String()

	// Serialize the paymentDetails struct to JSON
	paymentDetailsJSON, err := json.Marshal(paymentDetails)
	if err != nil {
		return "", err
	}

	query := `
        INSERT INTO transactions (transaction_id, status, payment_details)
        VALUES ($1, $2, $3)
    `
	_, err = DB.Exec(query, transactionID, status, paymentDetailsJSON)
	if err != nil {
		return "", err
	}
	return transactionID, nil
}

// UpdateTransactionStatus updates the status of a transaction.
func UpdateTransactionStatus(transactionID string, status string) error {
	query := `
        UPDATE transactions
        SET status = $1
        WHERE transaction_id = $2
    `
	_, err := DB.Exec(query, status, transactionID)
	return err
}
