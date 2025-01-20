package database

//
//import (
//	"log"
//	"os"
//
//	"github.com/joho/godotenv"
//	"gorm.io/driver/postgres"
//	"gorm.io/gorm"
//	//"main/internal/models"
//)
//
//var DB *gorm.DB
//
//func InitDB() error {
//	// Load environment variables from .env file
//	if err := godotenv.Load(); err != nil {
//		log.Fatal("Error loading .env file")
//	}
//
//	// Get the database credentials from the environment variables
//	host := os.Getenv("DB_HOST")
//	user := os.Getenv("DB_USER")
//	password := os.Getenv("DB_PASSWORD")
//	dbname := os.Getenv("DB_NAME")
//	port := os.Getenv("DB_PORT")
//	sslmode := os.Getenv("DB_SSLMODE")
//
//	// Construct the DSN (Data Source Name) string
//	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslmode
//
//	var err error
//	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
//	if err != nil {
//		return err
//	}
//
//	// Auto migrate models
//	return DB.AutoMigrate(&models.User{})
//}
