package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"recipe-generator/internal/api/config"
	"recipe-generator/internal/api/router"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	log.Printf("Starting Recipe Generator application...")

	// this is where my configuration will go
	log.Printf("Loading application configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("An error was encountered loading application config: %v\n", err)
	}
	log.Printf("Configuration loaded successfully. Environment: %s", cfg.Environment)

	// then I need to load the database connection.
	log.Printf("Establishing database connection...")
	connectionPool, err := loadDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to load connection to database: \n%v", err)
	}
	log.Printf("Database connection established successfully")

	// close the database when this function is exited
	defer func() {
		log.Printf("Closing database connection...")
		connectionPool.Close()
	}()

	// load any other connections I'll need such as any text message handling or router handling

	// load router file
	log.Printf("Initializing application router...")
	appRouter := router.NewRouter(cfg, connectionPool)
	log.Printf("Router initialized successfully")

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      appRouter.ServeMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start server using go routine. make sure you have TLS on for production.
	go func() {
		log.Printf("Starting HTTP server on port %s...", cfg.Port)
		var err error

		if cfg.Environment == "production" {
			log.Printf("Running in production mode with TLS enabled")
			err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			log.Printf("Running in %s mode without TLS", cfg.Environment)
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// wait for interrupt signal??? TODO, LEARN.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sig := <-quit
	log.Printf("Received shutdown signal: %v", sig)

	log.Println("Recipe Generator server is shutting down...")

	// create context for server shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// if I get an error shutting down the server, probably because it was being forced, then log.
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
	}

	log.Printf("Server shutdown complete")
}

// loads the database into the application using a connection string from the config.
func loadDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	log.Printf("Configuring database connection pool...")

	poolConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Set pool configuration
	poolConfig.MaxConns = cfg.MaxConns
	poolConfig.MaxConnLifetime = cfg.MaxConnLifetime
	poolConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	log.Printf("Database pool configuration: MaxConns=%d, MaxConnLifetime=%v, MaxConnIdleTime=%v",
		cfg.MaxConns, cfg.MaxConnLifetime, cfg.MaxConnIdleTime)

	// Create the connection pool
	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, err
	}

	// Test the connection
	if err := pool.Ping(context.Background()); err != nil {
		return nil, err
	}

	return pool, nil
}
