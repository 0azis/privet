package repository

import (
	"privet-api/internal/core/models"

	"github.com/jmoiron/sqlx"
)

type WishesRepository interface {
	Insert(wishes models.WishesList) error
	Select(userID string) (models.WishesList, error)
}

type wishes struct {
	db *sqlx.DB
}

func (w *wishes) Insert(wishes models.WishesList) error {
	if wishes.Gender == models.Female {
		wishes.Gender = models.Male
	} else {
		wishes.Gender = models.Female
	}	

	_, err := w.db.Exec(`insert into wishes (user_id, age, gender, location, activities) values ($1, $2, $3, $4, $5) on conflict (user_id) do update set age = $2, gender = $3, location = $4, activities = $5`, wishes.ID, wishes.Age, wishes.Gender, wishes.Location, wishes.Activities)
	return err
}

func (w *wishes) Select(userID string) (models.WishesList, error) {
	var result models.WishesList
	err := w.db.Get(&result, `select * from wishes where user_id = $1`, userID)
	return result, err
}

func NewWishesRepository(db *sqlx.DB) *wishes {
	return &wishes{
		db: db,
	}
}
