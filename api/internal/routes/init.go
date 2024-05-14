package routes

import (
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func InitRoutes(db *sqlx.DB) chi.Router {
	r := chi.NewRouter()
	api := r.Route("/v1", func(r chi.Router) {})

	GetUserRoutes(api, db)
	GetWishesRoutes(api, db)
	return r 
}
