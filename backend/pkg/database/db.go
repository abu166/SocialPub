package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// InitDB initializes the PostgreSQL database connection and performs auto-migration
func InitDB() error {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: No .env file found, using system environment variables")
	}

	// Construct DSN (Database Source Name)
	dsn := "host=" + os.Getenv("DB_HOST") +
		" user=" + os.Getenv("DB_USER") +
		" password=" + os.Getenv("DB_PASSWORD") +
		" dbname=" + os.Getenv("DB_NAME") +
		" port=" + os.Getenv("DB_PORT") +
		" sslmode=" + os.Getenv("DB_SSLMODE")

	// Open database connection
	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to the database:", err)
		return err
	}
	log.Println("Database connection established successfully!")

	// Test the connection with a simple query
	err = DB.Exec("SELECT 1").Error
	if err != nil {
		log.Fatal("Database connection test failed:", err)
		return err
	}
	log.Println("Database connection test passed!")

	// Auto-migrate tables
	if err := DB.AutoMigrate(&User{}, &Cart{}, &Transaction{}, &Product{}); err != nil {
		log.Fatal("Failed to auto-migrate tables:", err)
		return err
	}
	log.Println("Database tables migrated successfully!")

	// Reset the sequence for the carts table
	resetSequenceQuery := `
        SELECT setval('carts_cart_id_seq', COALESCE((SELECT MAX(cart_id) FROM carts), 0) + 1);
    `
	if err := DB.Exec(resetSequenceQuery).Error; err != nil {
		log.Printf("Failed to reset carts_cart_id_seq: %v", err)
		return err
	}
	log.Println("Reset carts_cart_id_seq successfully!")

	return nil
}
