package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/F1sssss/Perfect_Trade/internal/shared/config"
	"github.com/F1sssss/Perfect_Trade/internal/shared/logger"
)

// Server represents an HTTP server
type Server struct {
	httpServer *http.Server
	logger     logger.Logger
	config     *config.ServerConfig
}

// NewServer creates a new HTTP server
func NewServer(handler http.Handler, cfg *config.ServerConfig, log logger.Logger) *Server {
	return &Server{
		httpServer: &http.Server{
			Handler:      handler,
			ReadTimeout:  cfg.ReadTimeout,
			WriteTimeout: cfg.WriteTimeout,
			IdleTimeout:  cfg.IdleTimeout,
		},
		logger: log,
		config: cfg,
	}
}

// Start starts the HTTP server with graceful shutdown
func (s *Server) Start(port int) error {
	s.httpServer.Addr = fmt.Sprintf(":%d", port)

	// Channel to listen for errors
	serverErrors := make(chan error, 1)

	// Start server in a goroutine
	go func() {
		s.logger.Info("starting HTTP server", logger.Int("port", port))
		serverErrors <- s.httpServer.ListenAndServe()
	}()

	// Channel to listen for interrupt signals
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Block until error or shutdown signal
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
	case sig := <-shutdown:
		s.logger.Info("shutdown signal received", logger.String("signal", sig.String()))

		// Create context with timeout for shutdown
		ctx, cancel := context.WithTimeout(context.Background(), s.config.ShutdownTimeout)
		defer cancel()

		// Attempt graceful shutdown
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.Error("graceful shutdown failed", logger.Error(err))
			// Force close
			if closeErr := s.httpServer.Close(); closeErr != nil {
				return fmt.Errorf("force close error: %w", closeErr)
			}
			return fmt.Errorf("graceful shutdown error: %w", err)
		}

		s.logger.Info("server stopped gracefully")
		return nil
	}
}
