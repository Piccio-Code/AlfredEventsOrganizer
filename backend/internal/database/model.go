package database

import (
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	NotFoundError = errors.New("NOT FOUND")
)

type EventStatus string

const (
	EventStatusPending EventStatus = "pending"
	EventStatusCreated EventStatus = "created"
	EventStatusFailed  EventStatus = "failed"
)

type Models struct {
	// Add your domain models here, e.g.:
	// ExampleModel ExampleModel
	GroupModel        GroupModel
	PollTemplateModel PollTemplateModel
	PollModel         PollModel
}

func NewModels(DB *pgxpool.Pool) Models {
	return Models{
		GroupModel:        GroupModel{DB},
		PollTemplateModel: PollTemplateModel{DB},
		PollModel:         PollModel{DB},
	}
}
