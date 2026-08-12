package database

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	NotFoundError = errors.New("NOT FOUND")
)

type Models struct {
	// Add your domain models here, e.g.:
	// ExampleModel ExampleModel
	DB *pgxpool.Pool
}

func NewModels(DB *pgxpool.Pool) Models {
	return Models{
		DB: DB,
	}
}
