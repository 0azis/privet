package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/lib/pq"
)

type User struct {
	ID          int            `json:"id" db:"id"`
	Name        string         `json:"fullName" db:"name"`
	Age         int            `json:"age" db:"age"`
	Location    string         `json:"location" db:"location"`
	Gender      string         `json:"gender" db:"gender"`
	Description string         `json:"description" db:"description"`
	Activities  pq.StringArray `json:"activities" db:"activities"`
	// LastViewed int `json:"lastViewed" db:"last_viewed"`
}

const Female = "Женский"
const Male = "Мужской"

func (u *User) Validate() error {
	if u.Age < 14 {
		return errors.New("small age")
	}
	if u.ID == 0 || u.Name == "" || u.Age == 0 || u.Location == "" || u.Gender == "" {
		return errors.New("invalid data")
	}
	return nil
}

type httpError struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}

func JSONError(w http.ResponseWriter, code int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	he := httpError{
		StatusCode: code,
		Message:    message,
		Data:       data,
	}
	w.WriteHeader(he.StatusCode)
	json.NewEncoder(w).Encode(he)
}

type WishesList struct {
	ID         int            `json:"ID" db:"user_id"`
	Age        int            `json:"age" db:"age"`
	Gender     string         `json:"gender" db:"gender"`
	Location   string         `json:"city" db:"location"`
	Activities pq.StringArray `json:"activities" db:"activities"`
}

func (wl WishesList) Value() (driver.Value, error) {
	return json.Marshal(wl)
}
