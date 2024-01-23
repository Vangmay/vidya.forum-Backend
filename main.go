package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"github.com/vangmay/cvwo-backend/database"
	"github.com/vangmay/cvwo-backend/routes"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	} // Load environment variables to connect to database
	config := &database.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
	} // Create database config
	database.NewConnection(config) // Migrates the databse

	app := fiber.New()
	app.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowOrigins:     "http://localhost:8080",
		AllowHeaders:     "Origin, Content-Type, Accept, Accept-Language, Content-Length",
		AllowMethods:     "GET, POST, PATCH, DELETE, OPTIONS, PUT",
	}))

	routes.Setup(app) // Creation of routes

	app.Listen(":8080") // Server starts to listen
}
