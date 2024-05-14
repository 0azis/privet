package controllers

import (
	"fmt"
	"log/slog"
	"net/http"
	"privet-api/internal/core/models"
	"privet-api/internal/core/pkg"
	"privet-api/internal/store/repository"
)

type WishesControllers struct {
	repo repository.WishesRepository
}

func (wc *WishesControllers) AddWishes(w http.ResponseWriter, r *http.Request) {
	body, err := pkg.DecoderJSON("wishes", r.Body)
	fmt.Println(body)
	if err != nil {
		slog.Error("%s", err)
		models.JSONError(w,
			http.StatusBadRequest,
			http.StatusText(400),
			nil)
		return
	}

	wishes := body.(models.WishesList)
	err = wc.repo.Insert(wishes)
	if err != nil {
		slog.Error("%s", err)
		models.JSONError(w,
			http.StatusInternalServerError,
			http.StatusText(500),
			nil)
		return
	}
	models.JSONError(w,
		http.StatusCreated,
		http.StatusText(201),
		nil)
}

func (wc *WishesControllers) GetWishes(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		models.JSONError(w,
			http.StatusBadRequest,
			http.StatusText(400),
			nil)
		return
	}

	wishes, err := wc.repo.Select(id)
	if err != nil {
		slog.Error("%s", err)
		models.JSONError(w,
			http.StatusInternalServerError,
			http.StatusText(500),
			nil)
		return
	}

	models.JSONError(w,
		http.StatusOK,
		http.StatusText(200),
		wishes)
}

func NewWishesControllers(repo repository.WishesRepository) *WishesControllers {
	return &WishesControllers{
		repo,
	}
}
