package database

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"time"

	"github.com/Piccio-Code/AlfredEventsOranizer/backend/internal/waha"
)

type PollModel struct {
	DB *pgxpool.Pool
}

type PollRequest struct {
	TemplateID     string              `json:"template_id,omitempty"`
	Name           string              `json:"name,omitempty"`
	Title          string              `json:"title,omitempty"`
	GroupID        string              `json:"groupId,omitempty"`
	WhatsappChatId string              `json:"whatsappChatId,omitempty"`
	MultipleChoice bool                `json:"multipleChoice,omitempty"`
	Options        []PollOptionRequest `json:"options,omitempty"`
	ExpiresAt      time.Time           `json:"expiresAt,omitempty"`
}

type PollOptionRequest struct {
	Label                string `json:"label,omitempty"`
	Position             int    `json:"position,omitempty"`
	MotivationNeeded     bool   `json:"motivationNeeded,omitempty"`
	CongratulationNeeded bool   `json:"congratulationNeeded,omitempty"`
	SpecificationNeeded  bool   `json:"specificationNeeded,omitempty"`
}

func (r PollRequest) ToWahaPoll() waha.NewPollRequest {
	options := make([]string, 0, len(r.Options))
	for _, option := range r.Options {
		options = append(options, option.Label)
	}

	return waha.NewPollRequest{
		ChatId: r.WhatsappChatId,
		Poll: waha.PollPayload{
			Name:            r.Title,
			Options:         options,
			MultipleAnswers: r.MultipleChoice,
		},
	}
}

type PollResponse struct {
	ID             string               `json:"id"`
	TemplateID     string               `json:"templateId,omitempty"`
	WhatsappPollId string               `json:"whatsappPollId"`
	Name           string               `json:"name"`
	Title          string               `json:"title"`
	GroupID        string               `json:"groupId"`
	MultipleChoice bool                 `json:"multipleChoice"`
	Options        []PollOptionResponse `json:"options"`
	CreatedAt      time.Time            `json:"createdAt"`
	ExpiresAt      time.Time            `json:"expiresAt"`
	Status         string               `json:"status"`
}

type PollOptionResponse struct {
	ID                   string `json:"id"`
	PollID               string `json:"pollID"`
	Label                string `json:"label"`
	Position             int    `json:"position"`
	MotivationNeeded     bool   `json:"motivationNeeded"`
	CongratulationNeeded bool   `json:"congratulationNeeded"`
	SpecificationNeeded  bool   `json:"specificationNeeded"`
}

type PollDB struct {
	ID             string
	TemplateID     string
	WhatsAppPollID string
	Title          string
	MultipleChoice bool
	CreatedAt      time.Time
	ExpiresAt      time.Time
	GroupID        string
	Status         string
}

func (m *PollModel) Create(ctx context.Context, pollRequest PollRequest) (newPoll PollResponse, err error) {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := m.DB.Begin(ctx)

	defer tx.Rollback(ctx)

	if err != nil {
		return PollResponse{}, err
	}

	stmt := `
			INSERT INTO polls(template_id, title, multiple_choice, expires_at, group_id, status)
			VALUES (NULLIF($1, '')::uuid, $2, $3, $4, $5, $6)
			RETURNING id, template_id, title, multiple_choice, created_at, expires_at, group_id, status
			`

	args := []interface{}{pollRequest.TemplateID, pollRequest.Title, pollRequest.MultipleChoice, pollRequest.ExpiresAt, pollRequest.GroupID, EventStatusPending}

	err = tx.QueryRow(ctx, stmt, args...).Scan(&newPoll.ID, &newPoll.TemplateID, &newPoll.Title, &newPoll.MultipleChoice, &newPoll.CreatedAt, &newPoll.ExpiresAt, &newPoll.GroupID, &newPoll.Status)

	if err != nil {
		return PollResponse{}, err
	}

	var pollOptions []PollOptionResponse

	stmt = `
			INSERT INTO poll_options(poll_id, label, position, motivation_needed, congratulation_needed, specification_needed) 
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, poll_id, label, position, motivation_needed, congratulation_needed, specification_needed
			`

	for _, option := range pollRequest.Options {
		var newOption PollOptionResponse

		err := tx.QueryRow(ctx, stmt, newPoll.ID, option.Label, option.Position, option.MotivationNeeded, option.CongratulationNeeded, option.SpecificationNeeded).Scan(&newOption.ID, &newOption.PollID, &newOption.Label, &newOption.Position, &newOption.MotivationNeeded, &newOption.CongratulationNeeded, &newOption.SpecificationNeeded)

		if err != nil {
			return PollResponse{}, err
		}

		pollOptions = append(pollOptions, newOption)

	}

	return newPoll, tx.Commit(ctx)
}

func (m *PollModel) UpdateStatus(ctx context.Context, pollID, whatsappPollID string) error {
	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	stmt := `
			UPDATE polls
			SET whatsapp_poll_id = $1, status = $2
			WHERE id = $3
			`

	_, err := m.DB.Exec(ctx, stmt, whatsappPollID, EventStatusCreated, pollID)

	return err
}
