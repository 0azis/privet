package controllers

import (
	"log/slog"
	"net/http"
	"privet-api/internal/core/models"
	"privet-api/internal/core/pkg"
	"privet-api/internal/store/repository"
)

type UserControllers struct {
	wishesRepo repository.WishesRepository
	usersRepo  repository.UserRepository
}

func (uc *UserControllers) CreateAccount(w http.ResponseWriter, r *http.Request) {
	body, err := pkg.DecoderJSON("users", r.Body)
	if err != nil {
		slog.Error("%s", err)
		models.JSONError(w,
			http.StatusBadRequest,
			http.StatusText(400),
			nil)
		return
	}

	user := body.(models.User)
	if err := user.Validate(); err != nil {
		slog.Error("%v", err)
		models.JSONError(w,
			http.StatusBadRequest,
			http.StatusText(400),
			nil)
		return
	}

	wish := models.WishesList{
		ID:         user.ID,
		Age:        user.Age,
		Gender:     user.Gender,
		Location:   user.Location,
		Activities: user.Activities,
	}
	err = uc.usersRepo.Insert(wish, user)
	if err != nil {
		slog.Error("%v", err)
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
	slog.Info("user created...")
}

func (uc *UserControllers) GetUsers(w http.ResponseWriter, r *http.Request) {
	var sortedUsers []models.User // the final list with users

	id := r.URL.Query().Get("id")

	wishes, err := uc.wishesRepo.Select(id)
	if err != nil {
		slog.Error("%v", err)
		models.JSONError(w, http.StatusInternalServerError, http.StatusText(500), nil)
		return
	}

	resultUsers, err := uc.usersRepo.Select(id, wishes.Location, wishes.Gender)
	if err != nil {
		slog.Error("%v", err)
		models.JSONError(w, http.StatusInternalServerError, http.StatusText(500), nil)
		return
	}

	for x := range resultUsers {
		user := resultUsers[x]
		// check for possible age diffrence
		if user.Age-wishes.Age == 2 || user.Age-wishes.Age == -2 {
			sortedUsers = append(sortedUsers, user)
		}
		// check for the same activities
		for y := range user.Activities {
			for z := range wishes.Activities {
				if user.Activities[y] == wishes.Activities[z] {
					sortedUsers = append(sortedUsers, user)
					break
				}
			}
		}
	}

	if len(sortedUsers) == 0 {
		models.JSONError(w, http.StatusNotFound, http.StatusText(404), resultUsers)
		slog.Info("sending another users...")
		return
	}

	models.JSONError(w, http.StatusOK, http.StatusText(200), sortedUsers)
	slog.Info("sending users...")
}

func (uc *UserControllers) IsUserExist(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	isUser, _ := uc.usersRepo.GetOne(id)
	if isUser.ID == 0 {
		models.JSONError(w, http.StatusNotFound, http.StatusText(404), false)
		return
	} else {
		models.JSONError(w, http.StatusOK, http.StatusText(200), true)
		return
	}
}

func (uc *UserControllers) Profile(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")

	profile, err := uc.usersRepo.GetOne(id)
	if profile.ID == 0 {
		models.JSONError(w, http.StatusNotFound, http.StatusText(404), nil)
		return
	}

	if err != nil {
		models.JSONError(w, http.StatusInternalServerError, http.StatusText(500), nil)
		return
	}
	
	models.JSONError(w, http.StatusOK, http.StatusText(200), profile)
}

func NewUserControllers(wish repository.WishesRepository, user repository.UserRepository) *UserControllers {
	return &UserControllers{
		wishesRepo: wish,
		usersRepo:  user,
	}
}
