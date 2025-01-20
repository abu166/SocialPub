package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main/internal/models"
)

var DB *gorm.DB

func InitDB() error {
	dsn := "host=localhost user=abukhassymkhydyrbayev password=admin dbname=social_pub port=5432 sslmode=disable"

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	// Auto migrate models
	return DB.AutoMigrate(&models.User{})
}
