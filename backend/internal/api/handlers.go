package api

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HealthHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := contextWithTimeout(1 * time.Second)
		defer cancel()
		if err := db.PingContext(ctx); err != nil {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status": "unhealthy",
				"error":  err.Error(),
			})
		}
		return c.JSON(fiber.Map{"status": "ok"})
	}
}

func UsersHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		rows, err := db.Query("SELECT id, name, created_at FROM users ORDER BY id")
		if err != nil {
			log.Printf("query error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "query failed"})
		}
		defer rows.Close()

		var res []map[string]interface{}
		for rows.Next() {
			var id int
			var name string
			var created string
			if err := rows.Scan(&id, &name, &created); err != nil {
				log.Printf("row scan error: %v", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "row scan failed"})
			}
			res = append(res, fiber.Map{
				"id":         id,
				"name":       name,
				"created_at": created,
			})
		}
		if err := rows.Err(); err != nil {
			log.Printf("rows error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "rows iteration error"})
		}
		return c.JSON(res)
	}
}

// helper for creating a context with timeout
func contextWithTimeout(d time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), d)
}
