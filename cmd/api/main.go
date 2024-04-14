package main

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"gofiber-oauth/internal/server"
	"log"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	app := server.New()
	app.RegisterFiberRoutes()
	app.Listen(":8000")
}
