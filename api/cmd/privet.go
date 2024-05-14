package main

import (
	"log/slog"
	"net/http"
	"privet-api/internal/routes"
	"privet-api/internal/store"
)

func main() {
	slog.Info("http server started...")
	db, err := store.NewDB()
	if err != nil {
		slog.Error("database connection failed...")
	}

	r := routes.InitRoutes(db)

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		slog.Error("http server closed...")
	}

}
