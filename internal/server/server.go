package server

import (
	"github.com/gofiber/fiber/v2"

	"gofiber-oauth/internal/database"
)

type FiberServer struct {
	*fiber.App

	db database.Service
}

func New() *FiberServer {
	server := &FiberServer{
		App: fiber.New(fiber.Config{
			ServerHeader: "gofiber-oauth",
			AppName:      "gofiber-oauth",
		}),
	}
	return server
}
