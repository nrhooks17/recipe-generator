package router

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/config"
	"recipe-generator/internal/api/handlers"
)

// Router will need to add middleware later. Not now though.
type Router struct {
	ServeMux *http.ServeMux
}

func NewRouter(cfg *config.Config, db *pgxpool.Pool) *Router {

	router := &Router{}

	// init handlers here
	recipeHandler := handlers.NewRecipeHandler(db, cfg)
	healthHandler := handlers.HealthHandler{}

	mux := http.NewServeMux()

	// routes for application public routes can go here
	mux.Handle("/health", healthHandler.HealthCheck())

	mux.Handle("/", healthHandler.HealthCheck())

	mux.Handle("/recipe/random", recipeHandler.GetRandom())
	mux.Handle("/recipe/submit", recipeHandler.Post())
	mux.Handle("/recipe", recipeHandler.Get())

	// protected routes can go here.
	// r.Handle("/api/v1/user/profile", r.auth.Authenticate(userHandler.ProfileHandler()))
	router.ServeMux = mux

	return router
}
