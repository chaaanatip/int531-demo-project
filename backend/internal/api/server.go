package api

import (
	"database/sql"

	"github.com/chaaanatip/int531-demo-project/internal/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func NewApp(db *sql.DB) *fiber.App {
	app := fiber.New()

	// middleware

	app.Use(middleware.RequestIDMiddleware)
	app.Use(middleware.LoggerMiddleware)
	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://10.13.104.89",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Content-Type,Authorization",
	}))

	// health
	app.Get("/health", HealthHandler(db))

	// api group
	app.Get("/users", UsersHandler(db))

	return app
}
