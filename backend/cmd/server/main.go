package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/9inejames/int531-demo-project/internal/migration"

	"github.com/9inejames/int531-demo-project/internal/db"

	"github.com/9inejames/int531-demo-project/internal/config"

	"github.com/9inejames/int531-demo-project/internal/api"
)

func main() {
	// Load configuration
	cfg := config.LoadFromEnv()

	// Open DB (with retry)
	sqlDB, err := db.Open(cfg.DatabaseURL, db.Config{
		MaxOpenConns:    25,
		MaxIdleConns:    10,
		ConnMaxLifetime: 5 * time.Minute,
	})
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer sqlDB.Close()

	// Run migrations
	if err := migration.RunMigrations(sqlDB, cfg.MigrationsPath); err != nil {
		log.Fatalf("migrations failed: %v", err)
	}

	// Create fiber app and register routes
	app := api.NewApp(sqlDB)

	// Start server in goroutine
	serverErr := make(chan error, 1)
	go func() {
		log.Printf("listening on %s", cfg.ListenAddr)
		serverErr <- app.Listen(cfg.ListenAddr)
	}()

	// Graceful shutdown on SIGINT/SIGTERM
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-quit:
		log.Printf("signal received: %v — shutting down", sig)
	case err := <-serverErr:
		log.Printf("server error: %v", err)
	}

	// give app up to 10s to shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := app.Shutdown(); err != nil {
		log.Printf("app shutdown error: %v", err)
	}

	// ensure db closed (deferred)
	<-ctx.Done()
	log.Println("server stopped")
}
