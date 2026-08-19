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
	GroupModel        GroupModel
	PollTemplateModel PollTemplateModel
}

func NewModels(DB *pgxpool.Pool) Models {
	return Models{
		GroupModel:        GroupModel{DB},
		PollTemplateModel: PollTemplateModel{DB},
	}
}
