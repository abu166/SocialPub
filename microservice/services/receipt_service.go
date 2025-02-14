package services

import (
	"microservice/utils"
	"time"
)

func GenerateReceipt(transactionID string, details map[string]string) (string, error) {
	data := map[string]string{
		"TransactionID": transactionID,
		"OrderDate":     time.Now().Format("2006-01-02 15:04:05"),
		"ProductName":   "Subscription Plan",
		"PricePerUnit":  "$100",
		"Quantity":      "1",
		"TotalAmount":   "$100",
		"CustomerName":  details["name"],
		"PaymentMethod": "Credit Card (Encrypted)",
	}
	return utils.GeneratePDF(transactionID, data)
}
