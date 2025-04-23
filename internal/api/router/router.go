// Package router provides HTTP routing functionality for the recipe generator application.
// It defines routes and connects them to the appropriate handlers.
package router

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/config"
	"recipe-generator/internal/api/handler"
	"recipe-generator/internal/api/middleware"
)

// Router manages HTTP request routing for the application.
// It encapsulates an http.ServeMux and will support middleware in the future.
type Router struct {
	ServeMux http.Handler // The HTTP request multiplexer
}

// NewRouter creates and configures a new Router instance.
// It initializes all handlers, sets up routes, and returns the configured router.
// Parameters:
//   - cfg: Application configuration
//   - db: Database connection pool
//
// Returns a fully configured Router ready to handle HTTP requests.
func NewRouter(cfg *config.Config, db *pgxpool.Pool) *Router {

	router := &Router{}

	// init handlers here
	recipeHandler := handler.NewRecipeHandler(db, cfg)
	healthHandler := handler.HealthHandler{}

	mux := http.NewServeMux()

	// routes for application public routes can go here
	mux.Handle("/health", healthHandler.HealthCheck())

	mux.Handle("/", healthHandler.HealthCheck())

	mux.Handle("/recipe/random", recipeHandler.GetRandom())
	mux.Handle("/recipe/submit", recipeHandler.Post())
	mux.Handle("/recipe", recipeHandler.Get())

	// protected routes can go here.
	// r.Handle("/api/v1/user/profile", r.auth.Authenticate(userHandler.ProfileHandler()))

	
	router.ServeMux = middleware.CORS(mux) 

	return router
}
