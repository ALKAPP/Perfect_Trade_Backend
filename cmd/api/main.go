package main

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/F1sssss/Perfect_Trade/internal/shared/config"
	"github.com/F1sssss/Perfect_Trade/internal/shared/database"
	"github.com/F1sssss/Perfect_Trade/internal/shared/logger"
	"github.com/F1sssss/Perfect_Trade/internal/shared/server"
)

func main() {
	// Run application and exit with appropriate code
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	ctx := context.Background()

	// 1. Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// 2. Setup logger
	log, err := logger.NewLogger(&cfg.App)
	if err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}

	log.Info("starting application",
		logger.String("environment", cfg.App.Environment),
		logger.Int("port", cfg.App.Port),
	)

	// 3. Setup database
	pool, err := database.NewPostgresPool(ctx, &cfg.Database)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}
	defer database.Close(pool)

	log.Info("database connection established",
		logger.String("host", cfg.Database.Host),
		logger.Int("port", cfg.Database.Port),
		logger.String("database", cfg.Database.Name),
	)

	// 4. Setup router
	router := setupRouter()

	// 5. Setup and start HTTP server
	srv := server.NewServer(router, &cfg.Server, log)

	return srv.Start(cfg.App.Port)
}

func setupRouter() *chi.Mux {
	router := chi.NewRouter()

	// Middleware
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	// Health check endpoint
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"ok"}`))
	})

	// API routes (will add later)
	router.Route("/api/v1", func(r chi.Router) {
		// TODO: Register module routes here
		// r.Mount("/orders", orderHandler.Routes())
	})

	return router
}
