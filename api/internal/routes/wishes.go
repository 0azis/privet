package routes

import (
	"privet-api/internal/controllers"
	"privet-api/internal/store/repository"

	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func GetWishesRoutes(api chi.Router, db *sqlx.DB) {
	repo := repository.NewWishesRepository(db)
	controllers := controllers.NewWishesControllers(repo)

	api.Route("/wishes", func(r chi.Router) {
		r.Get("/", controllers.GetWishes)
		r.Post("/", controllers.AddWishes)
	})
}