package routes

import (
	"privet-api/internal/controllers"
	"privet-api/internal/store/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func GetUserRoutes(api chi.Router, db *sqlx.DB) {
	wishes := repository.NewWishesRepository(db)
	users := repository.NewUserRepository(db, wishes)
	controllers := controllers.NewUserControllers(wishes, users)

	api.Route("/users", func(r chi.Router) {
		r.Get("/", controllers.GetUsers)
		r.Get("/exist", controllers.IsUserExist)
		r.Get("/profile", controllers.Profile)
		r.Post("/", controllers.CreateAccount)
	})
}
