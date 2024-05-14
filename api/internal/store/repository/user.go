package repository

import (
	"privet-api/internal/core/models"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Insert(wish models.WishesList, user models.User) error
	Select(id string, location string, gener string) ([]models.User, error)
	GetOne(userID string) (models.User, error)
}

type user struct {
	wishes WishesRepository
	db     *sqlx.DB
}

func (u *user) Insert(wish models.WishesList, user models.User) error {
	_, err := u.db.Exec(`insert into users (id, name, age, location, gender, description, activities) values ($1, $2, $3, $4, $5, $6, $7)`, user.ID, user.Name, user.Age, user.Location, user.Gender, user.Description, user.Activities)
	if err != nil {
		return err
	}
	err = u.wishes.Insert(wish)
	if err != nil {
		return err
	}
	return nil
}

func (u *user) Select(id string, location string, gender string) ([]models.User, error) {
	users := []models.User{}
	err := u.db.Select(&users, `select * from users where id != $1 and lower(location) = lower($2) and gender = $3`, id, location, gender)
	return users, err
}

func (u *user) GetOne(userID string) (models.User, error) {
	var user models.User
	err := u.db.Get(&user, `select * from users where id = $1`, userID)
	return user, err
}

func NewUserRepository(db *sqlx.DB, wishes WishesRepository) *user {
	return &user{
		db:     db,
		wishes: wishes,
	}
}
