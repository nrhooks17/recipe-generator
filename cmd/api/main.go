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
	// this is where my configuration will go
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("An error was encountered loading application config: %v\n", err)
	}

	// then I need to load the database connection.
	connectionPool, err := loadDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to load connection to database: \n%v", err)
	}

	// close the database when this function is exited
	defer connectionPool.Close()

	// load any other connections I'll need such as any text message handling or router handling

	// load router file?
	router := router.NewRouter(cfg, connectionPool)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// start server using go routine. make sure you have TLS on for production.
	go func() {
		log.Printf("Starting server. Server is listening on port: %v", cfg.Port)
		var err error

		if cfg.Environment == "production" {
			err = server.ListenAndServeTLS(cfg.CertFile, cfg.KeyFile)
		} else {
			err = server.ListenAndServe()
		}

		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// wait for interrupt signal??? TODO, LEARN.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Recipe Generator server is shutting down...")

	// create context for server shutdown.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// if I get an error shutting down the server, probably because it was being forced, then log.
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server has been forced to shutdown: \n %v", err)
	}

	log.Println("Recipe generator has successfully shutdown properly.")
}

// loads the database into the application using a connection string from the config.
func loadDatabase(cfg *config.Config) (*pgxpool.Pool, error) {
	databaseConfig, err := pgxpool.ParseConfig(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// set database connection pool settings
	databaseConfig.MaxConns = cfg.MaxConns
	databaseConfig.MaxConnLifetime = cfg.MaxConnLifetime
	databaseConfig.MaxConnIdleTime = cfg.MaxConnIdleTime

	// connections should time out after 10 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	pool, err := pgxpool.NewWithConfig(ctx, databaseConfig)
	if err != nil {
		return nil, err
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, err
	}

	return pool, nil
}
