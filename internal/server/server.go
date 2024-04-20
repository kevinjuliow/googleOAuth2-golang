package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/session"

	"gofiber-oauth/internal/database"
)

type FiberServer struct {
	*fiber.App
	Store *session.Store
	db    database.Service
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
