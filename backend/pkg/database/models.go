package database

import (
	"encoding/json"
	"time"
)

// User represents the users table
type User struct {
	UserID    uint   `json:"user_id" gorm:"primaryKey"`
	UserName  string `json:"user_name" gorm:"not null;default:'Unknown'"`                     // Default value for UserName
	UserEmail string `json:"user_email" gorm:"not null;unique;default:'unknown@example.com'"` // Default value for UserEmail
	CreatedAt string `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
	UpdatedAt string `json:"updated_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Cart represents the structure of a cart in the database.
type Cart struct {
	CartID    uint      `json:"cart_id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id"`
	ProductID uint      `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Transaction represents the structure of a transaction in the database.
type Transaction struct {
	TransactionID  string          `json:"transaction_id" gorm:"primaryKey"`
	UserID         uint            `json:"user_id"`
	CartID         uint            `json:"cart_id"`
	Status         string          `json:"status"`
	PaymentDetails json.RawMessage `json:"payment_details"`
	CreatedAt      string          `json:"created_at" gorm:"default:CURRENT_TIMESTAMP"`
}

// Product represents the structure of a product in the database.
type Product struct {
	ProductID uint    `json:"product_id" gorm:"primaryKey"`
	Name      string  `json:"name" gorm:"not null"`
	Price     float64 `json:"price" gorm:"not null"`
}
