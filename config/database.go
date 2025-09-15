package config

import (
	"log"

	"blog-api-learn-go/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {
	dsn := "host=postgres user=postgres password=postgres dbname=blog port=5432 sslmode=disable"
	database, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database!", err)
	}
	DB = database
	// Auto migrate models
	DB.AutoMigrate(&models.User{}, &models.Post{}, &models.ActivityLog{})

	// Tạo GIN index cho cột tags
	DB.Exec("CREATE INDEX IF NOT EXISTS idx_posts_tags_gin ON posts USING GIN (tags);")
}
