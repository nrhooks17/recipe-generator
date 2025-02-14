package router

import (
	"encoding/json"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"recipe-generator/internal/api/config"
)

// will need to add middleware later. Not now though.
type Router struct {
	*http.ServeMux
}

func NewRouter(cfg *config.Config, db *pgxpool.Pool) *Router {
	router := &Router{
		ServeMux: http.NewServeMux(),
		// add middleware here
	}

	// add userHandler and authHandler here. Probably going to use custom uuid string first for authentication though.

	// routes for application public routes can go here
	router.Handle("/health-check", router.handleHealthCheck())

	// protected routes can go here.
	// r.Handle("/api/v1/user/profile", r.auth.Authenticate(userHandler.ProfileHandler()))

	return router
}

func (r *Router) handleHealthCheck() http.HandlerFunc {
	// return an anonymous function that returns a json with the status of my app

	return func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	}
}
