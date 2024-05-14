package store

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/jackc/pgx/stdlib"
)


func NewDB() (*sqlx.DB, error) {
	db, err := sqlx.Connect("pgx", "postgres://admin:admin123@localhost:5433/privet")
	return db, err	
} 
