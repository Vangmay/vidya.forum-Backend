package database

import (
	"fmt"

	"github.com/vangmay/cvwo-backend/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // Used to access the database

type Config struct {
	Host     string
	Port     string
	Password string
	User     string
	DBName   string
	SSLMode  string
}

// Create a DSN after extracting the environment variables
// Use gorm to open a postgressql database
// Error : Connection issues
// Automigrate the database
func NewConnection(config *Config) {
	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("could not connect to the database")
	}

	DB = db
	db.AutoMigrate(&models.User{}, &models.Post{}, &models.Comment{})

}
